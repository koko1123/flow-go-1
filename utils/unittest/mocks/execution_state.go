package mocks

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/engine/execution/state"

	"github.com/koko1123/flow-go-1/engine/execution/state/mock"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/storage"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

// ExecutionState is a mocked version of execution state that
// simulates some of its behavior for testing purpose
type ExecutionState struct {
	sync.Mutex
	mock.ExecutionState
	commits map[flow.Identifier]flow.StateCommitment
}

func NewExecutionState(seal *flow.Seal) *ExecutionState {
	commits := make(map[flow.Identifier]flow.StateCommitment)
	commits[seal.BlockID] = seal.FinalState
	return &ExecutionState{
		commits: commits,
	}
}

func (es *ExecutionState) PersistStateCommitment(ctx context.Context, blockID flow.Identifier, commit flow.StateCommitment) error {
	es.Lock()
	defer es.Unlock()
	es.commits[blockID] = commit
	return nil
}

func (es *ExecutionState) StateCommitmentByBlockID(ctx context.Context, blockID flow.Identifier) (flow.StateCommitment, error) {
	es.Lock()
	defer es.Unlock()
	commit, ok := es.commits[blockID]
	if !ok {
		return flow.DummyStateCommitment, storage.ErrNotFound
	}

	return commit, nil
}

func (es *ExecutionState) ExecuteBlock(t *testing.T, block *flow.Block) {
	parentExecuted, err := state.IsBlockExecuted(context.Background(), es, block.Header.ParentID)
	require.NoError(t, err)
	require.True(t, parentExecuted, "parent block not executed")
	require.NoError(t,
		es.PersistStateCommitment(
			context.Background(),
			block.ID(),
			unittest.StateCommitmentFixture()))
}
