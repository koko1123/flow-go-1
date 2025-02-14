package verification

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/consensus/hotstuff/mocks"
	"github.com/koko1123/flow-go-1/consensus/hotstuff/model"
	"github.com/koko1123/flow-go-1/consensus/hotstuff/signature"
	"github.com/onflow/flow-go/crypto"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module/local"
	modulemock "github.com/koko1123/flow-go-1/module/mock"
	msig "github.com/koko1123/flow-go-1/module/signature"
	"github.com/koko1123/flow-go-1/state/protocol"
	storagemock "github.com/koko1123/flow-go-1/storage/mock"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

// Test that when DKG key is available for a view, a signed block can pass the validation
// the sig include both staking sig and random beacon sig.
func TestCombinedSignWithDKGKey(t *testing.T) {
	identities := unittest.IdentityListFixture(4, unittest.WithRole(flow.RoleConsensus))

	// prepare data
	dkgKey := unittest.RandomBeaconPriv()
	pk := dkgKey.PublicKey()
	view := uint64(20)

	fblock := unittest.BlockFixture()
	fblock.Header.ProposerID = identities[0].NodeID
	fblock.Header.View = view
	block := model.BlockFromFlow(fblock.Header, 10)
	signerID := fblock.Header.ProposerID

	epochCounter := uint64(3)
	epochLookup := &modulemock.EpochLookup{}
	epochLookup.On("EpochForViewWithFallback", view).Return(epochCounter, nil)

	keys := &storagemock.SafeBeaconKeys{}
	// there is DKG key for this epoch
	keys.On("RetrieveMyBeaconPrivateKey", epochCounter).Return(dkgKey, true, nil)

	beaconKeyStore := signature.NewEpochAwareRandomBeaconKeyStore(epochLookup, keys)

	stakingPriv := unittest.StakingPrivKeyFixture()
	nodeID := unittest.IdentityFixture()
	nodeID.NodeID = signerID
	nodeID.StakingPubKey = stakingPriv.PublicKey()

	me, err := local.New(nodeID, stakingPriv)
	require.NoError(t, err)
	signer := NewCombinedSigner(me, beaconKeyStore)

	dkg := &mocks.DKG{}
	dkg.On("KeyShare", signerID).Return(pk, nil)

	committee := &mocks.Committee{}
	committee.On("DKG", mock.Anything).Return(dkg, nil)

	packer := signature.NewConsensusSigDataPacker(committee)
	verifier := NewCombinedVerifier(committee, packer)

	// check that a created proposal can be verified by a verifier
	proposal, err := signer.CreateProposal(block)
	require.NoError(t, err)

	vote := proposal.ProposerVote()
	err = verifier.VerifyVote(nodeID, vote.SigData, proposal.Block)
	require.NoError(t, err)

	// check that a created proposal's signature is a combined staking sig and random beacon sig
	msg := MakeVoteMessage(block.View, block.BlockID)
	stakingSig, err := stakingPriv.Sign(msg, msig.NewBLSHasher(msig.ConsensusVoteTag))
	require.NoError(t, err)

	beaconSig, err := dkgKey.Sign(msg, msig.NewBLSHasher(msig.RandomBeaconTag))
	require.NoError(t, err)

	expectedSig := msig.EncodeDoubleSig(stakingSig, beaconSig)
	require.Equal(t, expectedSig, proposal.SigData)

	// vote should be valid
	vote, err = signer.CreateVote(block)
	require.NoError(t, err)

	err = verifier.VerifyVote(nodeID, vote.SigData, block)
	require.NoError(t, err)

	// vote on different block should be invalid
	blockWrongID := *block
	blockWrongID.BlockID[0]++
	err = verifier.VerifyVote(nodeID, vote.SigData, &blockWrongID)
	require.ErrorIs(t, err, model.ErrInvalidSignature)

	// vote with a wrong view should be invalid
	blockWrongView := *block
	blockWrongView.View++
	err = verifier.VerifyVote(nodeID, vote.SigData, &blockWrongView)
	require.ErrorIs(t, err, model.ErrInvalidSignature)

	// vote by different signer should be invalid
	wrongVoter := identities[1]
	wrongVoter.StakingPubKey = unittest.StakingPrivKeyFixture().PublicKey()
	err = verifier.VerifyVote(wrongVoter, vote.SigData, block)
	require.ErrorIs(t, err, model.ErrInvalidSignature)

	// vote with changed signature should be invalid
	brokenSig := append([]byte{}, vote.SigData...) // copy
	brokenSig[4]++
	err = verifier.VerifyVote(nodeID, brokenSig, block)
	require.ErrorIs(t, err, model.ErrInvalidSignature)

	// Vote from a node that is _not_ part of the Random Beacon committee should be rejected.
	// Specifically, we expect that the verifier recognizes the `protocol.IdentityNotFoundError`
	// as a sign of an invalid vote and wraps it into a `model.InvalidSignerError`.
	*dkg = mocks.DKG{} // overwrite DKG mock with a new one
	dkg.On("KeyShare", signerID).Return(nil, protocol.IdentityNotFoundError{NodeID: signerID})
	err = verifier.VerifyVote(nodeID, vote.SigData, proposal.Block)
	require.True(t, model.IsInvalidSignerError(err))
}

// Test that when DKG key is not available for a view, a signed block can pass the validation
// the sig only include staking sig
func TestCombinedSignWithNoDKGKey(t *testing.T) {
	// prepare data
	dkgKey := unittest.RandomBeaconPriv()
	pk := dkgKey.PublicKey()
	view := uint64(20)

	fblock := unittest.BlockFixture()
	fblock.Header.View = view
	block := model.BlockFromFlow(fblock.Header, 10)
	signerID := fblock.Header.ProposerID

	epochCounter := uint64(3)
	epochLookup := &modulemock.EpochLookup{}
	epochLookup.On("EpochForViewWithFallback", view).Return(epochCounter, nil)

	keys := &storagemock.SafeBeaconKeys{}
	// there is no DKG key for this epoch
	keys.On("RetrieveMyBeaconPrivateKey", epochCounter).Return(nil, false, nil)

	beaconKeyStore := signature.NewEpochAwareRandomBeaconKeyStore(epochLookup, keys)

	stakingPriv := unittest.StakingPrivKeyFixture()
	nodeID := unittest.IdentityFixture()
	nodeID.NodeID = signerID
	nodeID.StakingPubKey = stakingPriv.PublicKey()

	me, err := local.New(nodeID, stakingPriv)
	require.NoError(t, err)
	signer := NewCombinedSigner(me, beaconKeyStore)

	dkg := &mocks.DKG{}
	dkg.On("KeyShare", signerID).Return(pk, nil)

	committee := &mocks.Committee{}
	// even if the node failed DKG, and has no random beacon private key,
	// but other nodes, who completed and succeeded DKG, have a public key
	// for this failed node, which can be used to verify signature from
	// this failed node.
	committee.On("DKG", mock.Anything).Return(dkg, nil)

	packer := signature.NewConsensusSigDataPacker(committee)
	verifier := NewCombinedVerifier(committee, packer)

	proposal, err := signer.CreateProposal(block)
	require.NoError(t, err)

	vote := proposal.ProposerVote()
	err = verifier.VerifyVote(nodeID, vote.SigData, proposal.Block)
	require.NoError(t, err)

	// As the proposer does not have a Random Beacon Key, it should sign solely with its staking key.
	// In this case, the SigData should be identical to the staking sig.
	expectedStakingSig, err := stakingPriv.Sign(
		MakeVoteMessage(block.View, block.BlockID),
		msig.NewBLSHasher(msig.ConsensusVoteTag),
	)
	require.NoError(t, err)
	require.Equal(t, expectedStakingSig, crypto.Signature(proposal.SigData))
}

// Test_VerifyQC_EmptySigners checks that Verifier returns an `model.InsufficientSignaturesError`
// if `signers` input is empty or nil. This check should happen _before_ the Verifier calls into
// any sub-components, because some (e.g. `crypto.AggregateBLSPublicKeys`) don't provide sufficient
// sentinel errors to distinguish between internal problems and external byzantine inputs.
func Test_VerifyQC_EmptySigners(t *testing.T) {
	committee := &mocks.Committee{}
	packer := signature.NewConsensusSigDataPacker(committee)
	verifier := NewCombinedVerifier(committee, packer)

	header := unittest.BlockHeaderFixture()
	block := model.BlockFromFlow(header, header.View-1)
	sigData := unittest.QCSigDataFixture()

	err := verifier.VerifyQC([]*flow.Identity{}, sigData, block)
	require.True(t, model.IsInsufficientSignaturesError(err))

	err = verifier.VerifyQC(nil, sigData, block)
	require.True(t, model.IsInsufficientSignaturesError(err))
}
