package finalizer

import (
	"github.com/koko1123/flow-go-1/consensus/hotstuff/model"
	"github.com/koko1123/flow-go-1/model/flow"
)

// BlockContainer wraps a block to implement forest.Vertex
// In addition, it holds some additional properties for efficient processing of blocks
// by the Finalizer
type BlockContainer struct {
	Block *model.Block
}

// functions implementing forest.vertex
func (b *BlockContainer) VertexID() flow.Identifier { return b.Block.BlockID }
func (b *BlockContainer) Level() uint64             { return b.Block.View }
func (b *BlockContainer) Parent() (flow.Identifier, uint64) {
	return b.Block.QC.BlockID, b.Block.QC.View
}
