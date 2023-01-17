//go:build relic
// +build relic

package verification

import (
	"errors"
	"fmt"

	"github.com/koko1123/flow-go-1/consensus/hotstuff"
	"github.com/koko1123/flow-go-1/consensus/hotstuff/model"
	"github.com/koko1123/flow-go-1/model/encoding"
	"github.com/koko1123/flow-go-1/model/flow"
	msig "github.com/koko1123/flow-go-1/module/signature"
	"github.com/koko1123/flow-go-1/state/protocol"
	"github.com/onflow/flow-go/crypto"
	"github.com/onflow/flow-go/crypto/hash"
)

// CombinedVerifierV3 is a verifier capable of verifying two signatures, one for each
// scheme. The first type is a signature from a staking signer,
// which verifies either a single or an aggregated signature. The second type is
// a signature from a random beacon signer, which verifies both the signature share and
// the reconstructed threshold signature.
type CombinedVerifierV3 struct {
	committee     hotstuff.Committee
	stakingHasher hash.Hasher
	beaconHasher  hash.Hasher
	packer        hotstuff.Packer
}

var _ hotstuff.Verifier = (*CombinedVerifierV3)(nil)

// NewCombinedVerifierV3 creates a new combined verifier with the given dependencies.
// - the hotstuff committee's state is used to retrieve the public keys for the staking signature;
// - the packer is used to unpack QC for verification;
func NewCombinedVerifierV3(committee hotstuff.Committee, packer hotstuff.Packer) *CombinedVerifierV3 {
	return &CombinedVerifierV3{
		committee:     committee,
		stakingHasher: msig.NewBLSHasher(msig.ConsensusVoteTag),
		beaconHasher:  msig.NewBLSHasher(msig.RandomBeaconTag),
		packer:        packer,
	}
}

// VerifyVote verifies the validity of a combined signature from a vote.
// Usually this method is only used to verify the proposer's vote, which is
// the vote included in a block proposal.
//   - model.InvalidFormatError if the signature has an incompatible format.
//   - model.ErrInvalidSignature is the signature is invalid
//   - model.InvalidSignerError if signer is _not_ part of the random beacon committee
//   - unexpected errors should be treated as symptoms of bugs or uncovered
//     edge cases in the logic (i.e. as fatal)
//
// This implementation already support the cases, where the DKG committee is a
// _strict subset_ of the full consensus committee.
func (c *CombinedVerifierV3) VerifyVote(signer *flow.Identity, sigData []byte, block *model.Block) error {

	// create the to-be-signed message
	msg := MakeVoteMessage(block.View, block.BlockID)

	sigType, sig, err := msig.DecodeSingleSig(sigData)
	if err != nil {
		if errors.Is(err, msig.ErrInvalidSignatureFormat) {
			return model.NewInvalidFormatErrorf("could not decode signature for block %v: %w", block.BlockID, err)
		}
		return fmt.Errorf("unexpected internal error while decoding signature for block %v: %w", block.BlockID, err)
	}

	switch sigType {
	case encoding.SigTypeStaking:
		// verify each signature against the message
		stakingValid, err := signer.StakingPubKey.Verify(sig, msg, c.stakingHasher)
		if err != nil {
			return fmt.Errorf("internal error while verifying staking signature for block %v: %w", block.BlockID, err)
		}
		if !stakingValid {
			return fmt.Errorf("invalid staking sig for block %v: %w", block.BlockID, model.ErrInvalidSignature)
		}

	case encoding.SigTypeRandomBeacon:
		dkg, err := c.committee.DKG(block.BlockID)
		if err != nil {
			return fmt.Errorf("could not get dkg: %w", err)
		}

		// if there is beacon share, there should be a beacon public key
		beaconPubKey, err := dkg.KeyShare(signer.NodeID)
		if err != nil {
			if protocol.IsIdentityNotFound(err) {
				return model.NewInvalidSignerErrorf("%v is not a random beacon participant: %w", signer.NodeID, err)
			}
			return fmt.Errorf("could not get random beacon key share for %x at block %v: %w", signer.NodeID, block.BlockID, err)
		}
		beaconValid, err := beaconPubKey.Verify(sig, msg, c.beaconHasher)
		if err != nil {
			return fmt.Errorf("internal error while verifying beacon signature for block %v: %w", block.BlockID, err)
		}
		if !beaconValid {
			return fmt.Errorf("invalid beacon sig for block %v: %w", block.BlockID, model.ErrInvalidSignature)
		}

	default:
		return model.NewInvalidFormatErrorf("invalid signature type %d", sigType)
	}

	return nil
}

// VerifyQC checks the cryptographic validity of the QC's `sigData` for the
// given block. It is the responsibility of the calling code to ensure
// that all `signers` are authorized, without duplicates. Return values:
//   - nil if `sigData` is cryptographically valid
//   - model.InsufficientSignaturesError if `signers` is empty.
//     Depending on the order of checks in the higher-level logic this error might
//     be an indicator of a external byzantine input or an internal bug.
//   - model.InvalidFormatError if `sigData` has an incompatible format
//   - model.ErrInvalidSignature if a signature is invalid
//   - model.InvalidSignerError if a signer is _not_ part of the random beacon committee
//   - error if running into any unexpected exception (i.e. fatal error)
//
// This implementation already support the cases, where the DKG committee is a
// _strict subset_ of the full consensus committee.
func (c *CombinedVerifierV3) VerifyQC(signers flow.IdentityList, sigData []byte, block *model.Block) error {
	if len(signers) == 0 {
		return model.NewInsufficientSignaturesErrorf("empty list of signers")
	}
	signerIdentities := signers.Lookup()
	dkg, err := c.committee.DKG(block.BlockID)
	if err != nil {
		return fmt.Errorf("could not get dkg data: %w", err)
	}

	// unpack sig data using packer
	blockSigData, err := c.packer.Unpack(signers, sigData)
	if err != nil {
		return fmt.Errorf("could not split signature: %w", err)
	}

	msg := MakeVoteMessage(block.View, block.BlockID)

	// STEP 1: verify random beacon group key
	// We do this first, since it is faster to check (no public key aggregation needed).
	beaconValid, err := dkg.GroupKey().Verify(blockSigData.ReconstructedRandomBeaconSig, msg, c.beaconHasher)
	if err != nil {
		return fmt.Errorf("internal error while verifying beacon signature: %w", err)
	}
	if !beaconValid {
		return fmt.Errorf("invalid reconstructed random beacon sig for block (%x): %w", block.BlockID, model.ErrInvalidSignature)
	}

	// verify the aggregated staking and beacon signatures next (more costly)
	// Caution: this function will error if pubKeys is empty
	verifyAggregatedSignature := func(pubKeys []crypto.PublicKey, aggregatedSig crypto.Signature, hasher hash.Hasher) error {
		// TODO: as further optimization, replace the following call with model/signature.PublicKeyAggregator
		aggregatedKey, err := crypto.AggregateBLSPublicKeys(pubKeys) // caution: requires non-empty slice of keys!
		if err != nil {
			return fmt.Errorf("internal error computing aggregated key: %w", err)
		}
		valid, err := aggregatedKey.Verify(aggregatedSig, msg, hasher)
		if err != nil {
			return fmt.Errorf("internal error while verifying aggregated signature: %w", err)
		}
		if !valid {
			return fmt.Errorf("invalid aggregated sig for block %v: %w", block.BlockID, model.ErrInvalidSignature)
		}
		return nil
	}

	// STEP 2: verify aggregated random beacon key shares
	// Step 2a: fetch all beacon signers public keys.
	// Note: A valid random beacon group sig is required for QC validity. To reconstruct
	// the group sig, _strictly more_ than `threshold` sig shares are required.
	threshold := msig.RandomBeaconThreshold(int(dkg.Size()))
	numRbSigners := len(blockSigData.RandomBeaconSigners)
	if numRbSigners <= threshold {
		// The Protocol prescribes that the random beacon signers that contributed to the QC are credited in the QC.
		// Depending on the reward model, under-reporting node contributions can be exploited in grieving attacks.
		// To construct a valid QC, the node generating it must have collected _more_ than `threshold` signatures.
		// Reporting fewer random beacon signers, the node is purposefully miss-representing node contributions.
		// We reject QCs with under-reported random beacon signers to reduce the surface of potential grieving attacks.
		return model.NewInvalidFormatErrorf("require at least %d random beacon sig shares but only got %d", threshold+1, numRbSigners)
	}
	beaconPubKeys := make([]crypto.PublicKey, 0, numRbSigners)
	for _, signerID := range blockSigData.RandomBeaconSigners {
		// Sanity check: every staking signer is in the list of authorized `signers`. (Thereby,
		// we enforce correctness within this component, as opposed relying on checks within the packer.)
		if _, ok := signerIdentities[signerID]; !ok {
			return fmt.Errorf("internal error, identity of random beacon signer not found %v", signerID)
		}
		keyShare, err := dkg.KeyShare(signerID)
		if err != nil {
			if protocol.IsIdentityNotFound(err) {
				return model.NewInvalidSignerErrorf("%v is not a random beacon participant: %w", signerID, err)
			}
			return fmt.Errorf("unexpected error retrieving dkg key share for signer %v: %w", signerID, err)
		}
		beaconPubKeys = append(beaconPubKeys, keyShare)
	}

	// Step 2b: verify aggregated beacon signature.
	// Our previous threshold check also guarantees that `beaconPubKeys` is not empty.
	err = verifyAggregatedSignature(beaconPubKeys, blockSigData.AggregatedRandomBeaconSig, c.beaconHasher)
	if err != nil {
		return fmt.Errorf("verifying aggregated random beacon sig shares failed for block %v: %w", block.BlockID, err)
	}

	// STEP 3: validating the aggregated staking signatures
	// Note: it is possible that all replicas signed with their random beacon keys, i.e.
	// `blockSigData.StakingSigners` could be empty. In this case, the
	// `blockSigData.AggregatedStakingSig` should also be empty.
	numStakingSigners := len(blockSigData.StakingSigners)
	if numStakingSigners == 0 {
		if len(blockSigData.AggregatedStakingSig) > 0 {
			return model.NewInvalidFormatErrorf("all replicas signed with random beacon keys, but QC has aggregated staking sig for block %v", block.BlockID)
		}
		// no aggregated staking sig to verify
		return nil
	}

	stakingPubKeys := make([]crypto.PublicKey, 0, numStakingSigners)
	for _, signerID := range blockSigData.StakingSigners {
		// Sanity check: every staking signer is in the list of authorized `signers`. (Thereby,
		// we enforce correctness within this component, as opposed relying on checks within the packer.)
		identity, ok := signerIdentities[signerID]
		if !ok {
			return fmt.Errorf("internal error, identity of staking signer not found %v", signerID)
		}
		stakingPubKeys = append(stakingPubKeys, identity.StakingPubKey)
	}
	err = verifyAggregatedSignature(stakingPubKeys, blockSigData.AggregatedStakingSig, c.stakingHasher)
	if err != nil {
		return fmt.Errorf("verifying aggregated staking sig failed for block %v: %w", block.BlockID, err)

	}

	return nil
}
