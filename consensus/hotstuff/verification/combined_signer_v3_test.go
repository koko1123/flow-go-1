package verification

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/consensus/hotstuff"
	"github.com/koko1123/flow-go-1/consensus/hotstuff/mocks"
	"github.com/koko1123/flow-go-1/consensus/hotstuff/model"
	"github.com/koko1123/flow-go-1/consensus/hotstuff/signature"
	"github.com/koko1123/flow-go-1/model/encoding"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module/local"
	modulemock "github.com/koko1123/flow-go-1/module/mock"
	msig "github.com/koko1123/flow-go-1/module/signature"
	"github.com/koko1123/flow-go-1/state/protocol"
	storagemock "github.com/koko1123/flow-go-1/storage/mock"
	"github.com/koko1123/flow-go-1/utils/unittest"
	"github.com/onflow/flow-go/crypto"
)

// Test that when DKG key is available for a view, a signed block can pass the validation
// the sig is a random beacon sig.
func TestCombinedSignWithDKGKeyV3(t *testing.T) {
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
	// there is DKG key for this epoch
	keys.On("RetrieveMyBeaconPrivateKey", epochCounter).Return(dkgKey, true, nil)

	beaconKeyStore := signature.NewEpochAwareRandomBeaconKeyStore(epochLookup, keys)

	stakingPriv := unittest.StakingPrivKeyFixture()
	nodeID := unittest.IdentityFixture()
	nodeID.NodeID = signerID
	nodeID.StakingPubKey = stakingPriv.PublicKey()

	me, err := local.New(nodeID, stakingPriv)
	require.NoError(t, err)
	signer := NewCombinedSignerV3(me, beaconKeyStore)

	dkg := &mocks.DKG{}
	dkg.On("KeyShare", signerID).Return(pk, nil)

	committee := &mocks.Committee{}
	committee.On("DKG", mock.Anything).Return(dkg, nil)

	packer := signature.NewConsensusSigDataPacker(committee)
	verifier := NewCombinedVerifierV3(committee, packer)

	// check that a created proposal can be verified by a verifier
	proposal, err := signer.CreateProposal(block)
	require.NoError(t, err)

	vote := proposal.ProposerVote()
	err = verifier.VerifyVote(nodeID, vote.SigData, proposal.Block)
	require.NoError(t, err)

	// check that a created proposal's signature is a combined staking sig and random beacon sig
	msg := MakeVoteMessage(block.View, block.BlockID)

	beaconSig, err := dkgKey.Sign(msg, msig.NewBLSHasher(msig.RandomBeaconTag))
	require.NoError(t, err)

	expectedSig := msig.EncodeSingleSig(encoding.SigTypeRandomBeacon, beaconSig)
	require.Equal(t, expectedSig, proposal.SigData)

	// Vote from a node that is _not_ part of the Random Beacon committee should be rejected.
	// Specifically, we expect that the verifier recognizes the `protocol.IdentityNotFoundError`
	// as a sign of an invalid vote and wraps it into a `model.InvalidSignerError`.
	*dkg = mocks.DKG{} // overwrite DKG mock with a new one
	dkg.On("KeyShare", signerID).Return(nil, protocol.IdentityNotFoundError{NodeID: signerID})
	err = verifier.VerifyVote(nodeID, vote.SigData, proposal.Block)
	require.True(t, model.IsInvalidSignerError(err))
}

// Test that when DKG key is not available for a view, a signed block can pass the validation
// the sig is a staking sig
func TestCombinedSignWithNoDKGKeyV3(t *testing.T) {
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
	signer := NewCombinedSignerV3(me, beaconKeyStore)

	dkg := &mocks.DKG{}
	dkg.On("KeyShare", signerID).Return(pk, nil)

	committee := &mocks.Committee{}
	// even if the node failed DKG, and has no random beacon private key,
	// but other nodes, who completed and succeeded DKG, have a public key
	// for this failed node, which can be used to verify signature from
	// this failed node.
	committee.On("DKG", mock.Anything).Return(dkg, nil)

	packer := signature.NewConsensusSigDataPacker(committee)
	verifier := NewCombinedVerifierV3(committee, packer)

	proposal, err := signer.CreateProposal(block)
	require.NoError(t, err)

	vote := proposal.ProposerVote()
	err = verifier.VerifyVote(nodeID, vote.SigData, proposal.Block)
	require.NoError(t, err)

	// check that a created proposal's signature is a combined staking sig and random beacon sig
	msg := MakeVoteMessage(block.View, block.BlockID)
	stakingSig, err := stakingPriv.Sign(msg, msig.NewBLSHasher(msig.ConsensusVoteTag))
	require.NoError(t, err)

	expectedSig := msig.EncodeSingleSig(encoding.SigTypeStaking, stakingSig)

	// check the signature only has staking sig
	require.Equal(t, expectedSig, proposal.SigData)
}

// Test_VerifyQC checks that a QC where either signer list is empty is rejected as invalid
func Test_VerifyQCV3(t *testing.T) {
	header := unittest.BlockHeaderFixture()
	block := model.BlockFromFlow(header, header.View-1)
	msg := MakeVoteMessage(block.View, block.BlockID)

	// generate some BLS key as a stub of the random beacon group key and use it to generate a reconstructed beacon sig
	privGroupKey, beaconSig := generateSignature(t, msg, msig.RandomBeaconTag)
	dkg := &mocks.DKG{}
	dkg.On("GroupKey").Return(privGroupKey.PublicKey(), nil)
	dkg.On("Size").Return(uint(20))
	committee := &mocks.Committee{}
	committee.On("DKG", mock.Anything).Return(dkg, nil)

	// generate 17 BLS keys as stubs for staking keys and use them to generate an aggregated staking sig
	privStakingKeys, aggStakingSig := generateAggregatedSignature(t, 17, msg, msig.ConsensusVoteTag)
	// generate 11 BLS keys as stubs for individual random beacon key shares and use them to generate an aggregated rand beacon sig
	privRbKeyShares, aggRbSig := generateAggregatedSignature(t, 11, msg, msig.RandomBeaconTag)

	stakingSigners := generateIdentitiesForPrivateKeys(t, privStakingKeys)
	rbSigners := generateIdentitiesForPrivateKeys(t, privRbKeyShares)
	registerPublicRbKeys(t, dkg, rbSigners.NodeIDs(), privRbKeyShares)
	allSigners := append(append(flow.IdentityList{}, stakingSigners...), rbSigners...)

	packedSigData := unittest.RandomBytes(1021)
	unpackedSigData := hotstuff.BlockSignatureData{
		StakingSigners:               stakingSigners.NodeIDs(),
		AggregatedStakingSig:         aggStakingSig,
		RandomBeaconSigners:          rbSigners.NodeIDs(),
		AggregatedRandomBeaconSig:    aggRbSig,
		ReconstructedRandomBeaconSig: beaconSig,
	}

	// first, we check that our testing setup works for a correct QC
	t.Run("valid QC", func(t *testing.T) {
		packer := &mocks.Packer{}
		packer.On("Unpack", mock.Anything, packedSigData).Return(&unpackedSigData, nil)

		verifier := NewCombinedVerifierV3(committee, packer)
		err := verifier.VerifyQC(allSigners, packedSigData, block)
		require.NoError(t, err)
	})

	// Here, we test correct verification of a QC, where all replicas signed with their
	// random beacon keys. This is optimal happy path.
	//  * empty list of staking signers
	//  * _no_ aggregated staking sig in QC
	// The Verifier should accept such QC
	t.Run("all replicas signed with random beacon keys", func(t *testing.T) {
		sd := unpackedSigData // copy correct QC
		sd.StakingSigners = []flow.Identifier{}
		sd.AggregatedStakingSig = []byte{}

		packer := &mocks.Packer{}
		packer.On("Unpack", mock.Anything, packedSigData).Return(&sd, nil)
		verifier := NewCombinedVerifierV3(committee, packer)
		err := verifier.VerifyQC(allSigners, packedSigData, block)
		require.NoError(t, err)
	})

	// Modify the correct QC:
	//  * empty list of staking signers
	//  * but an aggregated staking sig is given
	// The Verifier should recognize this as an invalid QC.
	t.Run("empty staking signers but aggregated staking sig in QC", func(t *testing.T) {
		sd := unpackedSigData // copy correct QC
		sd.StakingSigners = []flow.Identifier{}

		packer := &mocks.Packer{}
		packer.On("Unpack", mock.Anything, packedSigData).Return(&sd, nil)
		verifier := NewCombinedVerifierV3(committee, packer)
		err := verifier.VerifyQC(allSigners, packedSigData, block)
		require.True(t, model.IsInvalidFormatError(err))
	})

	// Modify the correct QC: empty list of random beacon signers.
	// The Verifier should recognize this as an invalid QC
	t.Run("empty random beacon signers", func(t *testing.T) {
		sd := unpackedSigData // copy correct QC
		sd.RandomBeaconSigners = []flow.Identifier{}

		packer := &mocks.Packer{}
		packer.On("Unpack", mock.Anything, packedSigData).Return(&sd, nil)
		verifier := NewCombinedVerifierV3(committee, packer)
		err := verifier.VerifyQC(allSigners, packedSigData, block)
		require.True(t, model.IsInvalidFormatError(err))
	})

	// Modify the correct QC: too few random beacon signers.
	// The Verifier should recognize this as an invalid QC
	t.Run("too few random beacon signers", func(t *testing.T) {
		// In total, we have 20 DKG participants, i.e. we require at least 10 random
		// beacon sig shares. But we only supply 5 aggregated key shares.
		sd := unpackedSigData // copy correct QC
		sd.RandomBeaconSigners = rbSigners[:5].NodeIDs()
		sd.AggregatedRandomBeaconSig = aggregatedSignature(t, privRbKeyShares[:5], msg, msig.RandomBeaconTag)

		packer := &mocks.Packer{}
		packer.On("Unpack", mock.Anything, packedSigData).Return(&sd, nil)
		verifier := NewCombinedVerifierV3(committee, packer)
		err := verifier.VerifyQC(allSigners, packedSigData, block)
		require.True(t, model.IsInvalidFormatError(err))
	})

}

// Test_VerifyQC_EmptySignersV3 checks that Verifier returns an `model.InsufficientSignaturesError`
// if `signers` input is empty or nil. This check should happen _before_ the Verifier calls into
// any sub-components, because some (e.g. `crypto.AggregateBLSPublicKeys`) don't provide sufficient
// sentinel errors to distinguish between internal problems and external byzantine inputs.
func Test_VerifyQC_EmptySignersV3(t *testing.T) {
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

func generateIdentitiesForPrivateKeys(t *testing.T, pivKeys []crypto.PrivateKey) flow.IdentityList {
	ids := make([]*flow.Identity, 0, len(pivKeys))
	for _, k := range pivKeys {
		id := unittest.IdentityFixture(
			unittest.WithRole(flow.RoleConsensus),
			unittest.WithStakingPubKey(k.PublicKey()),
		)
		ids = append(ids, id)
	}
	return ids
}

func registerPublicRbKeys(t *testing.T, dkg *mocks.DKG, signerIDs []flow.Identifier, pivKeys []crypto.PrivateKey) {
	assert.Equal(t, len(signerIDs), len(pivKeys), "one signer ID per key expected")
	for k, id := range signerIDs {
		dkg.On("KeyShare", id).Return(pivKeys[k].PublicKey(), nil)
	}
}

// generateAggregatedSignature generates `n` private BLS keys, signs `msg` which each key,
// and aggregates the resulting sigs. Returns private keys and aggregated sig.
func generateAggregatedSignature(t *testing.T, n int, msg []byte, tag string) ([]crypto.PrivateKey, crypto.Signature) {
	sigs := make([]crypto.Signature, 0, n)
	privs := make([]crypto.PrivateKey, 0, n)
	for ; n > 0; n-- {
		priv, sig := generateSignature(t, msg, tag)
		sigs = append(sigs, sig)
		privs = append(privs, priv)
	}
	agg, err := crypto.AggregateBLSSignatures(sigs)
	require.NoError(t, err)
	return privs, agg
}

// generateSignature creates a single private BLS 12-381 key, signs the provided `message` with
// using domain separation `tag` and return the private key and signature.
func generateSignature(t *testing.T, message []byte, tag string) (crypto.PrivateKey, crypto.Signature) {
	priv := unittest.PrivateKeyFixture(crypto.BLSBLS12381, crypto.KeyGenSeedMinLenBLSBLS12381)
	sig, err := priv.Sign(message, msig.NewBLSHasher(tag))
	require.NoError(t, err)
	return priv, sig
}

func aggregatedSignature(t *testing.T, pivKeys []crypto.PrivateKey, message []byte, tag string) crypto.Signature {
	hasher := msig.NewBLSHasher(tag)
	sigs := make([]crypto.Signature, 0, len(pivKeys))
	for _, k := range pivKeys {
		sig, err := k.Sign(message, hasher)
		require.NoError(t, err)
		sigs = append(sigs, sig)
	}
	agg, err := crypto.AggregateBLSSignatures(sigs)
	require.NoError(t, err)
	return agg
}
