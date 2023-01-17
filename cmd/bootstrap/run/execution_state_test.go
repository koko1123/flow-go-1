package run

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/fvm"
	"github.com/koko1123/flow-go-1/model/bootstrap"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

// This tests generates a checkpoint file to be used by the execution node when booting.
func TestGenerateExecutionState(t *testing.T) {
	seed := make([]byte, 48)
	seed[0] = 1
	sk, err := GenerateServiceAccountPrivateKey(seed)
	require.NoError(t, err)

	pk := sk.PublicKey(42)
	bootstrapDir := t.TempDir()
	trieDir := filepath.Join(bootstrapDir, bootstrap.DirnameExecutionState)
	commit, err := GenerateExecutionState(
		trieDir,
		pk,
		flow.Testnet.Chain(),
		fvm.WithInitialTokenSupply(unittest.GenesisTokenSupply))
	require.NoError(t, err)
	fmt.Printf("sk: %v\n", sk)
	fmt.Printf("pk: %v\n", pk)
	fmt.Printf("commit: %x\n", commit)
	fmt.Printf("a checkpoint file is generated at: %v\n", trieDir)
}
