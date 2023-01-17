package epochs

import (
	"context"
	"math/rand"
	"testing"

	"github.com/koko1123/flow-go-1/module"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	sdk "github.com/koko1123/flow-go-1-sdk"
	sdktemplates "github.com/koko1123/flow-go-1-sdk/templates"
	"github.com/koko1123/flow-go-1-sdk/test"

	hotstuff "github.com/koko1123/flow-go-1/consensus/hotstuff/mocks"
	hotstuffmodel "github.com/koko1123/flow-go-1/consensus/hotstuff/model"
	hotstuffver "github.com/koko1123/flow-go-1/consensus/hotstuff/verification"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module/epochs"
	modulemock "github.com/koko1123/flow-go-1/module/mock"
	"github.com/koko1123/flow-go-1/module/signature"
	"github.com/koko1123/flow-go-1/utils/unittest"

	clusterstate "github.com/koko1123/flow-go-1/state/cluster"
	protomock "github.com/koko1123/flow-go-1/state/protocol/mock"
)

func TestEmulatorBackedClusterQC(t *testing.T) {
	suite.Run(t, new(Suite))
}

// TestQuroumCertificate tests one Epoch of the EpochClusterQC contract
func (s *Suite) TestEpochQuorumCertificate() {

	// initial cluster and total node count
	epochCounter := uint64(1)
	clusterCount := 3
	nodesPerCluster := 10

	s.SetupTest()

	// create clustering with x clusters with x*y nodes
	clustering, nodes := s.CreateClusterList(clusterCount, nodesPerCluster)

	// mock the epoch object to return counter 0 and clustering as our clusterList
	epoch := &protomock.Epoch{}
	epoch.On("Counter").Return(epochCounter, nil)
	epoch.On("Clustering").Return(clustering, nil)

	s.PublishVoter()
	s.StartVoting(clustering, clusterCount, nodesPerCluster)

	// vote message to be signed with staking key
	blockID := unittest.IdentifierFixture()
	view := uint64(rand.Uint32())
	voteMessage := hotstuffver.MakeVoteMessage(view, blockID)

	// create cluster nodes with voter resource
	for _, node := range nodes {
		nodeID := node.NodeID
		stakingPrivKey := unittest.StakingPrivKeyFixture()

		// find cluster and create root block
		cluster, _, _ := clustering.ByNodeID(node.NodeID)
		rootBlock := clusterstate.CanonicalRootBlock(uint64(epochCounter), cluster)

		key, signer := test.AccountKeyGenerator().NewWithSigner()

		// create account on emualted chain
		address, err := s.blockchain.CreateAccount([]*sdk.AccountKey{key}, []sdktemplates.Contract{})
		s.Require().NoError(err)

		client := epochs.NewQCContractClient(zerolog.Nop(), s.emulatorClient, flow.ZeroID, nodeID, address.String(), 0, s.qcAddress.String(), signer)
		s.Require().NoError(err)

		local := &modulemock.Local{}
		local.On("NodeID").Return(nodeID)

		// create valid signature
		hasher := signature.NewBLSHasher(signature.CollectorVoteTag)
		signature, err := stakingPrivKey.Sign(voteMessage, hasher)
		s.Require().NoError(err)

		vote := hotstuffmodel.VoteFromFlow(nodeID, blockID, view, signature)
		hotSigner := &hotstuff.Signer{}
		hotSigner.On("CreateVote", mock.Anything).Return(vote, nil)

		snapshot := &protomock.Snapshot{}
		snapshot.On("Phase").Return(flow.EpochPhaseSetup, nil)

		state := &protomock.State{}
		state.On("CanonicalRootBlock").Return(rootBlock)
		state.On("Final").Return(snapshot)

		// create QC voter object to be used for voting for the root QC contract
		voter := epochs.NewRootQCVoter(zerolog.Logger{}, local, hotSigner, state, []module.QCContractClient{client})

		// create voter resource
		s.CreateVoterResource(address, nodeID, stakingPrivKey.PublicKey(), signer)

		// cast vote
		ctx := context.Background()
		err = voter.Vote(ctx, epoch)
		s.Require().NoError(err)
	}

	// check if each node has voted
	for _, node := range nodes {
		hasVoted := s.NodeHasVoted(node.NodeID)
		s.Assert().True(hasVoted)
	}

	// stop voting
	s.StopVoting()
}
