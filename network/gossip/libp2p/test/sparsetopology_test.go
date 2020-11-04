//nolint:unused
package test

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/onflow/flow-go/engine"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/model/libp2p/message"
	"github.com/onflow/flow-go/network/gossip/libp2p"
	"github.com/onflow/flow-go/network/gossip/libp2p/topology"
	"github.com/onflow/flow-go/utils/unittest"
)

// SparseTopologyTestSuite test 1-k messaging in a sparsely connected network
// Topology is used to control how the node is divided into subsets
type SparseTopologyTestSuite struct {
	suite.Suite
	ConduitWrapper
	nets []*libp2p.Network    // used to keep track of the networks
	mws  []*libp2p.Middleware // used to keep track of the middlewares associated with networks
	ids  flow.IdentityList    // used to keep track of the identifiers associated with networks
}

// TestSparseTopologyTestSuite runs all tests in this test suit
func TestSparseTopologyTestSuite(t *testing.T) {
	suite.Run(t, new(SparseTopologyTestSuite))
}

// TestSparselyConnectedNetworkScenario_Submit evaluates sparselyConnectedNetworkScenario on Submit method
// of Conduits
func (suite *SparseTopologyTestSuite) TestSparselyConnectedNetworkScenario_Submit() {
	suite.sparselyConnectedNetworkScenario(suite.Submit)
}

// TestSparselyConnectedNetworkScenario_Multicast evaluates sparselyConnectedNetworkScenario on Multicast
// method of Conduits
func (suite *SparseTopologyTestSuite) TestSparselyConnectedNetworkScenario_Multicast() {
	suite.sparselyConnectedNetworkScenario(suite.Multicast)
}

// TestSparselyConnectedNetworkScenario_Unicast evaluates sparselyConnectedNetworkScenario on Unicast
// method of Conduits
func (suite *SparseTopologyTestSuite) TestSparselyConnectedNetworkScenario_Unicast() {
	suite.sparselyConnectedNetworkScenario(suite.Unicast)
}

// TestSparselyConnectedNetworkScenario_Publish evaluates sparselyConnectedNetworkScenario on Publish
// method of Conduits
func (suite *SparseTopologyTestSuite) TestSparselyConnectedNetworkScenario_Publish() {
	suite.sparselyConnectedNetworkScenario(suite.Publish)
}

// sparselyConnectedNetworkScenario creates a network configuration with of 9 nodes with 3 subsets.
// The subsets are connected with each other with only one node
// 0,1,2,3 <-> 3,4,5,6 <-> 6,7,8,9
// Message sent by a node from one subset should be able to make it to nodes all subsets
func (suite *SparseTopologyTestSuite) sparselyConnectedNetworkScenario(send ConduitSendWrapperFunc) {
	// total number of nodes in the network
	const count = 9
	// total number of subnets (should be less than count)
	const subsets = 3

	ids, keys := GenerateIDs(suite.T(), count, RunNetwork)
	suite.ids = ids

	tops := createSparseTopology(count, subsets)

	// creates middleware and network instances
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
	mws := GenerateMiddlewares(suite.T(), logger, suite.ids, keys)
	sms := GenerateSubscriptionManagers(suite.T(), mws)
	suite.nets = GenerateNetworks(suite.T(), logger, suite.ids, mws, 100, tops, sms, RunNetwork)

	// create engines
	engs := make([]*MeshEngine, 0)
	for _, n := range suite.nets {
		eng := NewMeshEngine(suite.Suite.T(), n, count-1, engine.TestNetwork)
		engs = append(engs, eng)
	}

	// wait for nodes to heartbeat and discover each other
	time.Sleep(2 * time.Second)

	// node 0 broadcasting a message to all targets
	event := &message.TestMessage{
		Text: "hello from node 0",
	}
	require.NoError(suite.Suite.T(), send(event, engs[0].con, suite.ids.NodeIDs()...))

	// wait for message to be received by all recipients (excluding node 0)
	suite.checkMessageReception(engs, 1, count)

}

// TestDisjoinedNetworkScenario_Submit evaluates disjointedNetworkScenario using Submit method
// of conduits
func (suite *SparseTopologyTestSuite) TestDisjoinedNetworkScenario_Submit() {
	suite.disjointedNetworkScenario(suite.Submit)
}

// TestDisjoinedNetworkScenario_Publish evaluates disjointedNetworkScenario using Publish method
// of conduits
func (suite *SparseTopologyTestSuite) TestDisjoinedNetworkScenario_Publish() {
	suite.disjointedNetworkScenario(suite.Publish)
}

// TestDisjoinedNetworkScenario_Multicast evaluates disjointedNetworkScenario using Multicast method
// of conduits
func (suite *SparseTopologyTestSuite) TestDisjoinedNetworkScenario_Multicast() {
	suite.disjointedNetworkScenario(suite.Multicast)
}

// disjointedNetworkScenario creates a network configuration of 9 nodes with 3 subsets.
// The subsets are not connected with each other.
// 0,1,2 | 3,4,5 | 6,7,8
// Message sent by a node from one subset should not be able to make to nodes in a different subset
// This test is created primarily to prove that topology is indeed honored by the networking layer since technically,
// each node does have the ip addresses of all other nodes and could just disregard topology all together and connect
// to every other node directly making the sparselyConnectedNetworkScenario test meaningless
func (suite *SparseTopologyTestSuite) disjointedNetworkScenario(send ConduitSendWrapperFunc) {
	// total number of nodes in the network
	const count = 9
	// total number of subnets (should be less than count)
	const subnets = 3

	ids, keys := GenerateIDs(suite.T(), count, RunNetwork)
	suite.ids = ids

	tops := createDisjointedTopology(count, subnets)

	// creates middleware and network instances
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
	mws := GenerateMiddlewares(suite.T(), logger, suite.ids, keys)
	sms := GenerateSubscriptionManagers(suite.T(), mws)
	suite.nets = GenerateNetworks(suite.T(), logger, suite.ids, mws, 100, tops, sms, RunNetwork)

	// create engines
	engs := make([]*MeshEngine, 0)
	for _, n := range suite.nets {
		eng := NewMeshEngine(suite.Suite.T(), n, count-1, engine.TestNetwork)
		engs = append(engs, eng)
	}

	// wait for nodes to heartbeat and discover each other
	// this is a sparse network so it may need a at least 3 seconds (1 for each subnet)
	time.Sleep(4 * time.Second)

	// node 0 broadcasting a message to ALL targets
	event := &message.TestMessage{
		Text: "hello from node 0",
	}
	require.NoError(suite.Suite.T(), send(event, engs[0].con, suite.ids.NodeIDs()...))

	// wait for message to be received by nodes only in subset 1 (excluding node 0)
	suite.checkMessageReception(engs, 1, subnets)
}

// TearDownTest closes the networks within a specified timeout
func (suite *SparseTopologyTestSuite) TearDownTest() {
	stopNetworks(suite.T(), suite.nets, 3*time.Second)
	suite.ids = nil
	suite.mws = nil
	suite.nets = nil
}

// IndexBoundTopology is a topology implementation that limits the subset by indices of the identity list
type IndexBoundTopology struct {
	minIndex int
	maxIndex int
}

// Returns a subset of ids bounded by [minIndex, maxIndex) for the SparseTopology
func (ibt IndexBoundTopology) Subset(idList flow.IdentityList, _ uint, _ string) (flow.IdentityList, error) {
	sub := idList[ibt.minIndex:ibt.maxIndex]
	return sub, nil
}

// createSparseTopology creates topologies for nodes such that subsets have one overlapping node
// e.g. top 1 - 0,1,2,3; top 2 - 3,4,5,6; top 3 - 6,7,8,9
func createSparseTopology(count int, subsets int) []topology.Topology {
	tops := make([]topology.Topology, count)
	subsetLen := count / subsets
	for i := 0; i < count; i++ {
		s := i / subsets          // which subset does this node belong to
		minIndex := s * subsetLen //minIndex is just a multiple subset length
		var maxIndex int
		if s == subsets-1 { // if this is the last subset
			maxIndex = count // then set max index to count since nodes so may not evenly divide by # of subsets
		} else {
			maxIndex = ((s + 1) * subsetLen) + 1 // a plus one to cause an overlap between subsets
		}
		var top topology.Topology = IndexBoundTopology{
			minIndex: minIndex,
			maxIndex: maxIndex,
		}
		tops[i] = top
	}
	return tops
}

// createDisjointedTopology creates topologies for nodes such that subsets don't have any overlap
// e.g. top 1 - 0,1,2; top 2 - 3,4,5; top 3 - 6,7,8
func createDisjointedTopology(count int, subsets int) []topology.Topology {
	tops := make([]topology.Topology, count)
	subsetLen := count / subsets
	for i := 0; i < count; i++ {
		s := i / subsets          // which subset does this node belong to
		minIndex := s * subsetLen //minIndex is just a multiple subset length
		var maxIndex int
		if s == subsets-1 { // if this is the last subset
			maxIndex = count // then set max index to count since nodes so may not evenly divide by # of subsets
		} else {
			maxIndex = (s + 1) * subsetLen
		}

		var top topology.Topology = IndexBoundTopology{
			minIndex: minIndex,
			maxIndex: maxIndex,
		}

		tops[i] = top
	}
	return tops
}

// checkMessageReception checks if engs[low:high) have received a message while all the other engs have not
func (suite SparseTopologyTestSuite) checkMessageReception(engs []*MeshEngine, low int, high int) {
	wg := sync.WaitGroup{}
	// fires a goroutine for all engines to listens for the incoming message
	for _, e := range engs[low:high] {
		wg.Add(1)
		go func(e *MeshEngine) {
			<-e.received
			wg.Done()
		}(e)
	}

	unittest.RequireReturnsBefore(suite.T(), wg.Wait, 10*time.Second, "test timed out on broadcast dissemination")

	// evaluates that all messages are received
	for i, e := range engs {
		if i >= low && i < high {
			assert.Len(suite.Suite.T(), e.event, 1, fmt.Sprintf("engine %d did not receive the message", i))
		} else {
			assert.Len(suite.Suite.T(), e.event, 0, fmt.Sprintf("engine %d received the message", i))
		}
	}
}
