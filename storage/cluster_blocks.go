package storage

import (
	"github.com/koko1123/flow-go-1/model/cluster"
	"github.com/koko1123/flow-go-1/model/flow"
)

type ClusterBlocks interface {

	// Store stores the cluster block.
	Store(block *cluster.Block) error

	// ByID returns the block with the given ID.
	ByID(blockID flow.Identifier) (*cluster.Block, error)

	// ByHeight returns the block with the given height. Only available for
	// finalized blocks.
	ByHeight(height uint64) (*cluster.Block, error)
}
