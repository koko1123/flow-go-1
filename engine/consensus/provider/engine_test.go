// (c) 2019 Dapper Labs - ALL RIGHTS RESERVED

package provider

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module/metrics"
	module "github.com/koko1123/flow-go-1/module/mock"
	"github.com/koko1123/flow-go-1/module/trace"
	"github.com/koko1123/flow-go-1/network/mocknetwork"
	protocol "github.com/koko1123/flow-go-1/state/protocol/mock"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

type Suite struct {
	suite.Suite

	me      *module.Local
	conduit *mocknetwork.Conduit
	state   *protocol.State
	final   *protocol.Snapshot

	identities flow.IdentityList

	engine *Engine
}

func TestProviderEngine(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (suite *Suite) SetupTest() {

	suite.me = new(module.Local)
	suite.conduit = new(mocknetwork.Conduit)
	suite.state = new(protocol.State)
	suite.final = new(protocol.Snapshot)

	suite.engine = &Engine{
		me:      suite.me,
		state:   suite.state,
		con:     suite.conduit,
		message: metrics.NewNoopCollector(),
		tracer:  trace.NewNoopTracer(),
	}

	suite.identities = unittest.CompleteIdentitySet()
	localID := suite.identities[0].NodeID

	suite.me.On("NodeID").Return(localID)
	suite.state.On("Final").Return(suite.final)
	suite.final.On("Identities", mock.Anything).Return(
		func(f flow.IdentityFilter) flow.IdentityList { return suite.identities.Filter(f) },
		func(flow.IdentityFilter) error { return nil },
	)
}

// proposals submitted by remote nodes should not be accepted.
func (suite *Suite) TestOnBlockProposal_RemoteOrigin() {

	proposal := unittest.ProposalFixture()
	// message submitted by remote node
	err := suite.engine.onBlockProposal(suite.identities[1].NodeID, proposal)
	suite.Assert().Error(err)
}

func (suite *Suite) OnBlockProposal_Success() {

	proposal := unittest.ProposalFixture()

	params := []interface{}{proposal}
	for _, identity := range suite.identities {
		// skip consensus nodes
		if identity.Role == flow.RoleConsensus {
			continue
		}
		params = append(params, identity.NodeID)
	}

	suite.conduit.On("Publish", params...).Return(nil).Once()

	err := suite.engine.onBlockProposal(suite.me.NodeID(), proposal)
	suite.Require().Nil(err)
	suite.conduit.AssertExpectations(suite.T())
}
