package execution_data

import (
	"github.com/ipfs/go-cid"

	"github.com/koko1123/flow-go-1/ledger"
	"github.com/koko1123/flow-go-1/model/flow"
)

const DefaultMaxBlobSize = 1 << 20 // 1MiB

// ChunkExecutionData represents the execution data of a chunk
type ChunkExecutionData struct {
	Collection *flow.Collection
	Events     flow.EventsList
	TrieUpdate *ledger.TrieUpdate
}

type BlockExecutionDataRoot struct {
	BlockID               flow.Identifier
	ChunkExecutionDataIDs []cid.Cid
}

type BlockExecutionData struct {
	BlockID             flow.Identifier
	ChunkExecutionDatas []*ChunkExecutionData
}
