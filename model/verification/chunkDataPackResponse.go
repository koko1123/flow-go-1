package verification

import (
	"github.com/koko1123/flow-go-1/model/chunks"
	"github.com/koko1123/flow-go-1/model/flow"
)

// ChunkDataPackResponse is an internal data structure in fetcher engine that is passed between the fetcher
// and requester engine. It conveys requested chunk data pack as well as meta-data for fetcher engine to
// process the chunk data pack.
type ChunkDataPackResponse struct {
	chunks.Locator
	Cdp *flow.ChunkDataPack
}
