package jobqueue

import (
	"fmt"

	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module"
	"github.com/koko1123/flow-go-1/state/protocol"
	"github.com/koko1123/flow-go-1/storage"
)

// FinalizedBlockReader provides an abstraction for consumers to read blocks as job.
type FinalizedBlockReader struct {
	state  protocol.State
	blocks storage.Blocks
}

var _ module.Jobs = (*FinalizedBlockReader)(nil)

// NewFinalizedBlockReader creates and returns a FinalizedBlockReader.
func NewFinalizedBlockReader(state protocol.State, blocks storage.Blocks) *FinalizedBlockReader {
	return &FinalizedBlockReader{
		state:  state,
		blocks: blocks,
	}
}

// AtIndex returns the block job at the given index.
// The block job at an index is just the finalized block at that index (i.e., height).
func (r FinalizedBlockReader) AtIndex(index uint64) (module.Job, error) {
	block, err := r.blockByHeight(index)
	if err != nil {
		return nil, fmt.Errorf("could not get block by index %v: %w", index, err)
	}
	return BlockToJob(block), nil
}

// blockByHeight returns the block at the given height.
func (r FinalizedBlockReader) blockByHeight(height uint64) (*flow.Block, error) {
	block, err := r.blocks.ByHeight(height)
	if err != nil {
		return nil, fmt.Errorf("could not get block by height %d: %w", height, err)
	}

	return block, nil
}

// Head returns the last finalized height as job index.
func (r FinalizedBlockReader) Head() (uint64, error) {
	header, err := r.state.Final().Head()
	if err != nil {
		return 0, fmt.Errorf("could not get header of last finalized block: %w", err)
	}

	return header.Height, nil
}
