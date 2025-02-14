package jobqueue_test

import (
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/engine/testutil"
	vertestutils "github.com/koko1123/flow-go-1/engine/verification/utils/unittest"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/model/flow/filter"
	"github.com/koko1123/flow-go-1/module/jobqueue"
	"github.com/koko1123/flow-go-1/module/metrics"
	"github.com/koko1123/flow-go-1/module/trace"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

// TestBlockReader evaluates that block reader correctly reads stored finalized blocks from the blocks storage and
// protocol state.
func TestBlockReader(t *testing.T) {
	withReader(t, 10, func(reader *jobqueue.FinalizedBlockReader, blocks []*flow.Block) {
		// head of block reader should be the same height as the last block on the chain.
		head, err := reader.Head()
		require.NoError(t, err)
		require.Equal(t, head, blocks[len(blocks)-1].Header.Height)

		// retrieved blocks from block reader should be the same as the original blocks stored in it.
		for _, actual := range blocks {
			index := actual.Header.Height
			job, err := reader.AtIndex(index)
			require.NoError(t, err)

			retrieved, err := jobqueue.JobToBlock(job)
			require.NoError(t, err)
			require.Equal(t, actual.ID(), retrieved.ID())
		}
	})
}

// withReader is a test helper that sets up a block reader.
// It also provides a chain of specified number of finalized blocks ready to read by block reader, i.e., the protocol state is extended with the
// chain of blocks and the blocks are stored in blocks storage.
func withReader(
	t *testing.T,
	blockCount int,
	withBlockReader func(*jobqueue.FinalizedBlockReader, []*flow.Block),
) {
	require.Equal(t, blockCount%2, 0, "block count for this test should be even")
	unittest.RunWithBadgerDB(t, func(db *badger.DB) {

		collector := &metrics.NoopCollector{}
		tracer := trace.NewNoopTracer()
		participants := unittest.IdentityListFixture(5, unittest.WithAllRoles())
		rootSnapshot := unittest.RootSnapshotFixture(participants)
		s := testutil.CompleteStateFixture(t, collector, tracer, rootSnapshot)

		reader := jobqueue.NewFinalizedBlockReader(s.State, s.Storage.Blocks)

		// generates a chain of blocks in the form of root <- R1 <- C1 <- R2 <- C2 <- ... where Rs are distinct reference
		// blocks (i.e., containing guarantees), and Cs are container blocks for their preceding reference block,
		// Container blocks only contain receipts of their preceding reference blocks. But they do not
		// hold any guarantees.
		root, err := s.State.Params().Root()
		require.NoError(t, err)
		clusterCommittee := participants.Filter(filter.HasRole(flow.RoleCollection))
		results := vertestutils.CompleteExecutionReceiptChainFixture(t, root, blockCount/2, vertestutils.WithClusterCommittee(clusterCommittee))
		blocks := vertestutils.ExtendStateWithFinalizedBlocks(t, results, s.State)

		withBlockReader(reader, blocks)
	})
}
