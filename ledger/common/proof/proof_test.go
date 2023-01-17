package proof_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/ledger/common/proof"
	"github.com/koko1123/flow-go-1/ledger/common/testutils"
)

// Test_ProofVerify tests proof verification
func Test_TrieProofVerify(t *testing.T) {
	p, sc := testutils.TrieProofFixture()
	require.True(t, proof.VerifyTrieProof(p, sc))
}

// Test_BatchProofVerify tests batch proof verification
func Test_TrieBatchProofVerify(t *testing.T) {
	bp, sc := testutils.TrieBatchProofFixture()
	require.True(t, proof.VerifyTrieBatchProof(bp, sc))
}
