package signature

import (
	"bytes"
	"crypto/rand"
	"errors"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/consensus/hotstuff"
	"github.com/onflow/flow-go/crypto"
	_ "github.com/onflow/flow-go/engine"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/utils/unittest"
)

func sortIdentities(ids []flow.Identity) []flow.Identity {
	canonicalOrder := func(i, j int) bool {
		return bytes.Compare(ids[i].NodeID[:], ids[j].NodeID[:]) < 0
	}
	sort.Slice(ids, canonicalOrder)
	return ids
}

func sortIdentifiers(ids []flow.Identifier) []flow.Identifier {
	canonicalOrder := func(i, j int) bool {
		return bytes.Compare(ids[i][:], ids[j][:]) < 0
	}
	sort.Slice(ids, canonicalOrder)
	return ids
}

func createAggregationData(t *testing.T, signersNumber int) (
	hotstuff.WeightedSignatureAggregator, []flow.Identity, []crypto.Signature, []byte, string) {
	// create identities
	ids := make([]flow.Identity, 0, signersNumber)
	for i := 0; i < signersNumber; i++ {
		ids = append(ids, *unittest.IdentityFixture())
	}
	ids = sortIdentities(ids)

	// create message and tag
	msgLen := 100
	msg := make([]byte, msgLen)
	tag := "random_tag"
	hasher := crypto.NewBLSKMAC(tag)

	// create keys, identities and signatures
	sigs := make([]crypto.Signature, 0, signersNumber)
	seed := make([]byte, crypto.KeyGenSeedMinLenBLSBLS12381)
	for i := 0; i < signersNumber; i++ {
		// keys
		_, err := rand.Read(seed)
		require.NoError(t, err)
		sk, err := crypto.GeneratePrivateKey(crypto.BLSBLS12381, seed)
		require.NoError(t, err)
		ids[i].StakingPubKey = sk.PublicKey()
		// signatures
		sig, err := sk.Sign(msg, hasher)
		require.NoError(t, err)
		sigs = append(sigs, sig)
	}
	aggregator, err := NewWeightedSignatureAggregator(ids, msg, tag)
	require.NoError(t, err)
	return aggregator, ids, sigs, msg, tag
}

func verifyAggregate(signers []flow.Identifier, ids []flow.Identity, sig []byte,
	msg []byte, tag string) (bool, error) {
	// query identity using identifier
	// done linearly just for testing
	getIdentity := func(signer flow.Identifier) *flow.Identity {
		for _, id := range ids {
			if id.NodeID == signer {
				return &id
			}
		}
		return nil
	}

	// get keys
	keys := make([]crypto.PublicKey, 0, len(ids))
	for _, signer := range signers {
		id := getIdentity(signer)
		if id == nil {
			return false, errors.New("unexpected test error")
		}
		keys = append(keys, id.StakingPubKey)
	}
	// verify signature
	hasher := crypto.NewBLSKMAC(tag)
	return crypto.VerifyBLSSignatureOneMessage(keys, sig, msg, hasher)
}

func TestWeightedSignatureAggregator(t *testing.T) {
	signersNum := 20

	// constrcutor edge cases
	t.Run("constructor", func(t *testing.T) {
		msg := []byte("random_msg")
		tag := "random_tag"

		signer := unittest.IdentityFixture()
		// identity with empty key
		_, err := NewWeightedSignatureAggregator([]flow.Identity{*signer}, msg, tag)
		assert.Error(t, err)
		// wrong key types
		seed := make([]byte, crypto.KeyGenSeedMinLenECDSAP256)
		_, err = rand.Read(seed)
		require.NoError(t, err)
		sk, err := crypto.GeneratePrivateKey(crypto.ECDSAP256, seed)
		require.NoError(t, err)
		signer.StakingPubKey = sk.PublicKey()
		_, err = NewWeightedSignatureAggregator([]flow.Identity{*signer}, msg, tag)
		assert.Error(t, err)
		// empty signers
		_, err = NewWeightedSignatureAggregator([]flow.Identity{}, msg, tag)
		assert.Error(t, err)
	})

	// Happy paths
	t.Run("happy path", func(t *testing.T) {
		aggregator, ids, sigs, msg, tag := createAggregationData(t, signersNum)
		// only add half of the signatures
		subSet := signersNum / 2
		var expectedWeight uint64
		for i, sig := range sigs[subSet:] {
			index := i + subSet
			// test Verify
			err := aggregator.Verify(ids[index].NodeID, sig)
			assert.NoError(t, err)
			// test TrustedAdd
			weight, err := aggregator.TrustedAdd(ids[index].NodeID, sig)
			assert.NoError(t, err)
			expectedWeight += ids[index].Stake
			assert.Equal(t, expectedWeight, weight)
		}
		signers, agg, err := aggregator.Aggregate()
		assert.NoError(t, err)
		ok, err := verifyAggregate(signers, ids, agg, msg, tag)
		assert.NoError(t, err)
		assert.True(t, ok)
		// check signers
		signers = sortIdentifiers(signers)
		for i := 0; i < subSet; i++ {
			index := i + subSet
			assert.Equal(t, signers[i], ids[index].NodeID)
		}

		// add remaining signatures
		for i, sig := range sigs[:subSet] {
			weight, err := aggregator.TrustedAdd(ids[i].NodeID, sig)
			assert.NoError(t, err)
			expectedWeight += ids[i].Stake
			assert.Equal(t, expectedWeight, weight)
		}
		signers, agg, err = aggregator.Aggregate()
		assert.NoError(t, err)
		ok, err = verifyAggregate(signers, ids, agg, msg, tag)
		assert.NoError(t, err)
		assert.True(t, ok)
		// check signers
		signers = sortIdentifiers(signers)
		for i := 0; i < signersNum; i++ {
			assert.Equal(t, signers[i], ids[i].NodeID)
		}
	})
	/*
		invalidInput := engine.NewInvalidInputError("some error")
		duplicate := newErrDuplicatedSigner("some error")

		// Unhappy paths
		t.Run("invalid inputs", func(t *testing.T) {
			aggregator, sigs := createAggregationData(t, signersNum)
			// loop through invalid inputs
			for _, index := range []int{-1, signersNum} {
				ok, err := aggregator.Verify(index, sigs[0])
				assert.False(t, ok)
				assert.Error(t, err)
				assert.IsType(t, invalidInput, err)

				ok, err = aggregator.VerifyAndAdd(index, sigs[0])
				assert.False(t, ok)
				assert.Error(t, err)
				assert.IsType(t, invalidInput, err)

				err = aggregator.TrustedAdd(index, sigs[0])
				assert.Error(t, err)
				assert.IsType(t, invalidInput, err)

				ok, err = aggregator.HasSignature(index)
				assert.False(t, ok)
				assert.Error(t, err)
				assert.IsType(t, invalidInput, err)

				ok, err = aggregator.VerifyAggregate([]int{index}, sigs[0])
				assert.False(t, ok)
				assert.Error(t, err)
				assert.IsType(t, invalidInput, err)
			}
			// empty list
			ok, err := aggregator.VerifyAggregate([]int{}, sigs[0])
			assert.False(t, ok)
			assert.Error(t, err)
			assert.IsType(t, invalidInput, err)
		})

		t.Run("duplicate signature", func(t *testing.T) {
			aggregator, sigs := createAggregationData(t, signersNum)
			for i, sig := range sigs {
				err := aggregator.TrustedAdd(i, sig)
				require.NoError(t, err)
			}
			// TrustedAdd
			for i := range sigs {
				err := aggregator.TrustedAdd(i, sigs[i]) // same signature for same index
				assert.Error(t, err)
				assert.IsType(t, duplicate, err)
				err = aggregator.TrustedAdd(0, sigs[(i+1)%signersNum]) // different signature for same index
				assert.Error(t, err)
				assert.IsType(t, duplicate, err)
				// VerifyAndAdd
				ok, err := aggregator.VerifyAndAdd(i, sigs[i]) // valid but redundant signature
				assert.False(t, ok)
				assert.Error(t, err)
				assert.IsType(t, duplicate, err)
			}
		})

		t.Run("invalid signature", func(t *testing.T) {
			aggregator, sigs := createAggregationData(t, signersNum)
			// corrupt sigs[0]
			sigs[0][4] ^= 1
			// test Verify
			ok, err := aggregator.Verify(0, sigs[0])
			require.NoError(t, err)
			assert.False(t, ok)
			// test Verify and Add
			ok, err = aggregator.VerifyAndAdd(0, sigs[0])
			require.NoError(t, err)
			assert.False(t, ok)
			// check signature is still not added
			ok, err = aggregator.HasSignature(0)
			require.NoError(t, err)
			assert.False(t, ok)
			// add signatures for aggregation including corrupt sigs[0]
			for i, sig := range sigs {
				err := aggregator.TrustedAdd(i, sig)
				require.NoError(t, err)
			}
			signers, agg, err := aggregator.Aggregate()
			assert.Error(t, err)
			assert.Nil(t, agg)
			assert.Nil(t, signers)
			// fix sigs[0]
			sigs[0][4] ^= 1
		})

		// cached aggregated signature
		t.Run("cached aggregated signature", func(t *testing.T) {

		})*/
}
