// (c) 2019 Dapper Labs - ALL RIGHTS RESERVED

package ingestion

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/flow-go/model/flow"
	mockmodule "github.com/dapperlabs/flow-go/module/mock"
	mocknetwork "github.com/dapperlabs/flow-go/network/mock"
	mockprotocol "github.com/dapperlabs/flow-go/state/protocol/mock"
	"github.com/dapperlabs/flow-go/utils/unittest"
)

func TestOnCollectionGuaranteeValid(t *testing.T) {

	prop := &mocknetwork.Engine{}
	state := &mockprotocol.State{}
	final := &mockprotocol.Snapshot{}
	metrics := &mockmodule.Metrics{}
	metrics.On("StartCollectionToFinalized", mock.Anything).Return()

	e := &Engine{
		prop:    prop,
		state:   state,
		metrics: metrics,
	}

	originID := unittest.IdentifierFixture()
	guarantee := unittest.CollectionGuaranteeFixture()
	guarantee.SignerIDs = []flow.Identifier{originID}
	guarantee.Signature = unittest.SignatureFixture()

	identity := unittest.IdentityFixture(unittest.WithRole(flow.RoleCollection))
	identity.NodeID = originID
	clusters := flow.NewClusterList(1)
	clusters.Add(0, identity)

	header := unittest.BlockHeaderFixture()
	state.On("Final").Return(final).Once()
	state.On("AtBlockID", mock.Anything).Return(final).Once()
	final.On("Identity", originID).Return(identity, nil).Once()
	final.On("Head").Return(&header, nil).Twice()
	final.On("Clusters").Return(clusters, nil).Once()
	prop.On("SubmitLocal", guarantee).Return().Once()

	err := e.onCollectionGuarantee(originID, guarantee)
	require.NoError(t, err)

	state.AssertExpectations(t)
	final.AssertExpectations(t)
	prop.AssertExpectations(t)
}

func TestOnCollectionGuaranteeMissingIdentity(t *testing.T) {

	prop := &mocknetwork.Engine{}
	state := &mockprotocol.State{}
	final := &mockprotocol.Snapshot{}

	e := &Engine{
		prop:  prop,
		state: state,
	}

	originID := unittest.IdentifierFixture()
	guarantee := unittest.CollectionGuaranteeFixture()
	guarantee.SignerIDs = []flow.Identifier{originID}
	guarantee.Signature = unittest.SignatureFixture()

	identity := unittest.IdentityFixture(unittest.WithRole(flow.RoleCollection))
	identity.NodeID = originID
	clusters := flow.NewClusterList(1)
	clusters.Add(0, identity)

	state.On("Final").Return(final).Once()
	final.On("Identity", originID).Return(nil, errors.New("identity error")).Once()

	err := e.onCollectionGuarantee(originID, guarantee)
	require.Error(t, err)

	state.AssertExpectations(t)
	final.AssertExpectations(t)
	prop.AssertExpectations(t)
}

func TestOnCollectionGuaranteeInvalidRole(t *testing.T) {

	prop := &mocknetwork.Engine{}
	state := &mockprotocol.State{}
	final := &mockprotocol.Snapshot{}

	e := &Engine{
		prop:  prop,
		state: state,
	}

	originID := unittest.IdentifierFixture()
	guarantee := unittest.CollectionGuaranteeFixture()
	guarantee.SignerIDs = []flow.Identifier{originID}
	guarantee.Signature = unittest.SignatureFixture()

	// origin node has wrong role (consensus)
	identity := unittest.IdentityFixture(unittest.WithRole(flow.RoleConsensus))
	identity.NodeID = originID
	clusters := flow.NewClusterList(1)
	clusters.Add(0, identity)

	state.On("Final").Return(final).Once()
	final.On("Identity", originID).Return(identity, nil).Once()

	err := e.onCollectionGuarantee(originID, guarantee)
	require.Error(t, err)

	state.AssertExpectations(t)
	final.AssertExpectations(t)
	prop.AssertExpectations(t)
}

func TestOnCollectionGuaranteeExpired(t *testing.T) {

	prop := &mocknetwork.Engine{}
	state := &mockprotocol.State{}
	final := &mockprotocol.Snapshot{}
	past := &mockprotocol.Snapshot{}
	metrics := &mockmodule.Metrics{}
	metrics.On("StartCollectionToFinalized", mock.Anything).Return()

	e := &Engine{
		prop:    prop,
		state:   state,
		metrics: metrics,
	}

	originID := unittest.IdentifierFixture()
	guarantee := unittest.CollectionGuaranteeFixture()
	guarantee.SignerIDs = []flow.Identifier{originID}
	guarantee.Signature = unittest.SignatureFixture()

	identity := unittest.IdentityFixture(unittest.WithRole(flow.RoleCollection))
	identity.NodeID = originID
	clusters := flow.NewClusterList(1)
	clusters.Add(0, identity)

	finalBlk := unittest.BlockHeaderFixture()
	finalBlk.Height = flow.DefaultTransactionExpiry + 10 // head has moved 10 blocks beyond the transaction expiry limit

	refBlk := unittest.BlockHeaderFixture()
	refBlk.Height = 0
	guarantee.ReferenceBlockID = refBlk.ID() // guarantee points to a reference block in the past and is expired

	state.On("Final").Return(final).Once()
	state.On("AtBlockID", refBlk.ID()).Return(past).Once()

	final.On("Identity", originID).Return(identity, nil).Once()
	final.On("Head").Return(&finalBlk, nil).Once()

	past.On("Head").Return(&refBlk, nil).Once()

	err := e.onCollectionGuarantee(originID, guarantee)
	require.Error(t, err)

	state.AssertExpectations(t)
	final.AssertExpectations(t)
	prop.AssertExpectations(t)
}
