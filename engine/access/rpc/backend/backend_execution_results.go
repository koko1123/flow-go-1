package backend

import (
	"context"

	"github.com/koko1123/flow-go-1/engine/common/rpc"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/storage"
)

type backendExecutionResults struct {
	executionResults storage.ExecutionResults
}

func (b *backendExecutionResults) GetExecutionResultForBlockID(ctx context.Context, blockID flow.Identifier) (*flow.ExecutionResult, error) {
	er, err := b.executionResults.ByBlockID(blockID)
	if err != nil {
		return nil, rpc.ConvertStorageError(err)
	}

	return er, nil
}

// GetExecutionResultByID gets an execution result by its ID.
func (b *backendExecutionResults) GetExecutionResultByID(ctx context.Context, id flow.Identifier) (*flow.ExecutionResult, error) {
	result, err := b.executionResults.ByID(id)
	if err != nil {
		return nil, rpc.ConvertStorageError(err)
	}

	return result, nil
}
