// Package cluster contains models related to collection node cluster
// consensus.
package cluster

import (
	"github.com/koko1123/flow-go-1/model/flow"
)

func Genesis() *Block {
	header := &flow.Header{
		View:      0,
		ChainID:   "cluster",
		Timestamp: flow.GenesisTime,
		ParentID:  flow.ZeroID,
	}

	payload := EmptyPayload(flow.ZeroID)

	block := &Block{
		Header: header,
	}
	block.SetPayload(payload)

	return block
}

// Block represents a block in collection node cluster consensus. It contains
// a standard block header with a payload containing only a single collection.
type Block struct {
	Header  *flow.Header
	Payload *Payload
}

// ID returns the ID of the underlying block header.
func (b Block) ID() flow.Identifier {
	return b.Header.ID()
}

// SetPayload sets the payload and payload hash.
func (b *Block) SetPayload(payload Payload) {
	b.Payload = &payload
	b.Header.PayloadHash = payload.Hash()
}
