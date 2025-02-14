package chunks

import (
	"github.com/koko1123/flow-go-1/model/flow"
)

// ChunkListFromCommit creates a chunklist with one chunk whos final state is
// the commit
func ChunkListFromCommit(commit flow.StateCommitment) flow.ChunkList {
	chunks := flow.ChunkList{}
	chunk := &flow.Chunk{
		Index:    0,
		EndState: commit,
	}
	chunks.Insert(chunk)

	return chunks
}
