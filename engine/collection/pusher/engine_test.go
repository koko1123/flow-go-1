package pusher_test

import (
	"io"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/koko1123/flow-go-1/engine/collection/pusher"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/model/flow/filter"
	"github.com/koko1123/flow-go-1/model/messages"
	"github.com/koko1123/flow-go-1/module/metrics"
	module "github.com/koko1123/flow-go-1/module/mock"
	"github.com/koko1123/flow-go-1/network/channels"
	"github.com/koko1123/flow-go-1/network/mocknetwork"
	protocol "github.com/koko1123/flow-go-1/state/protocol/mock"
	storage "github.com/koko1123/flow-go-1/storage/mock"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

type Suite struct {
	suite.Suite

	identities   flow.IdentityList
	state        *protocol.State
	snapshot     *protocol.Snapshot
	conduit      *mocknetwork.Conduit
	me           *module.Local
	collections  *storage.Collections
	transactions *storage.Transactions

	engine *pusher.Engine
}

func (suite *Suite) SetupTest() {
	var err error

	// add some dummy identities so we have one of each role
	suite.identities = unittest.IdentityListFixture(5, unittest.WithAllRoles())
	me := suite.identities.Filter(filter.HasRole(flow.RoleCollection))[0]

	suite.state = new(protocol.State)
	suite.snapshot = new(protocol.Snapshot)
	suite.snapshot.On("Identities", mock.Anything).Return(func(filter flow.IdentityFilter) flow.IdentityList {
		return suite.identities.Filter(filter)
	}, func(filter flow.IdentityFilter) error {
		return nil
	})
	suite.state.On("Final").Return(suite.snapshot)

	metrics := metrics.NewNoopCollector()

	net := new(mocknetwork.Network)
	suite.conduit = new(mocknetwork.Conduit)
	net.On("Register", mock.Anything, mock.Anything).Return(suite.conduit, nil)

	suite.me = new(module.Local)
	suite.me.On("NodeID").Return(me.NodeID)

	suite.collections = new(storage.Collections)
	suite.transactions = new(storage.Transactions)

	suite.engine, err = pusher.New(
		zerolog.New(io.Discard),
		net,
		suite.state,
		metrics,
		metrics,
		suite.me,
		suite.collections,
		suite.transactions,
	)
	suite.Require().Nil(err)
}

func TestPusherEngine(t *testing.T) {
	suite.Run(t, new(Suite))
}

// should be able to submit collection guarantees to consensus nodes
func (suite *Suite) TestSubmitCollectionGuarantee() {

	guarantee := unittest.CollectionGuaranteeFixture()

	// should submit the collection to consensus nodes
	consensus := suite.identities.Filter(filter.HasRole(flow.RoleConsensus))
	suite.conduit.On("Publish", guarantee, consensus[0].NodeID).Return(nil)

	msg := &messages.SubmitCollectionGuarantee{
		Guarantee: *guarantee,
	}
	err := suite.engine.ProcessLocal(msg)
	suite.Require().Nil(err)

	suite.conduit.AssertExpectations(suite.T())
}

// should be able to submit collection guarantees to consensus nodes
func (suite *Suite) TestSubmitCollectionGuaranteeNonLocal() {

	guarantee := unittest.CollectionGuaranteeFixture()

	// send from a non-allowed role
	sender := suite.identities.Filter(filter.HasRole(flow.RoleVerification))[0]

	msg := &messages.SubmitCollectionGuarantee{
		Guarantee: *guarantee,
	}
	err := suite.engine.Process(channels.PushGuarantees, sender.NodeID, msg)
	suite.Require().Error(err)

	suite.conduit.AssertNumberOfCalls(suite.T(), "Multicast", 0)
}
