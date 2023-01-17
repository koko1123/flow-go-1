package badgermodel

import (
	"github.com/koko1123/flow-go-1/model/flow"
)

// StoredChunkDataPack is an in-storage representation of chunk data pack.
// Its prime difference is instead of an actual collection, it keeps a collection ID hence relying on maintaining
// the collection on a secondary storage.
type StoredChunkDataPack struct {
	ChunkID      flow.Identifier
	StartState   flow.StateCommitment
	Proof        flow.StorageProof
	CollectionID flow.Identifier
	SystemChunk  bool
}
