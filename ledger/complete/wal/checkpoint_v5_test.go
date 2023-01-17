package wal

import (
	"path/filepath"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/utils/unittest"
)

func TestCopyCheckpointFileV5(t *testing.T) {
	unittest.RunWithTempDir(t, func(dir string) {
		tries := createSimpleTrie(t)
		fileName := "checkpoint"
		logger := unittest.Logger()
		require.NoErrorf(t, StoreCheckpointV5(dir, fileName, &logger, tries...), "fail to store checkpoint")
		to := filepath.Join(dir, "newfolder")
		newPaths, err := CopyCheckpointFile(fileName, dir, to)
		require.NoError(t, err)
		log.Info().Msgf("copied to :%v", newPaths)
		decoded, err := LoadCheckpoint(filepath.Join(to, fileName), &logger)
		require.NoErrorf(t, err, "fail to read checkpoint %v/%v", dir, fileName)
		requireTriesEqual(t, tries, decoded)
	})
}
