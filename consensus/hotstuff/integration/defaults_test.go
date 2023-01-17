package integration

import (
	"time"

	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

func DefaultRoot() *flow.Header {
	header := &flow.Header{
		ChainID:     "chain",
		ParentID:    flow.ZeroID,
		Height:      0,
		PayloadHash: unittest.IdentifierFixture(),
		Timestamp:   time.Now().UTC(),
	}
	return header
}

func DefaultStart() uint64 {
	return 1
}

func DefaultPruned() uint64 {
	return 0
}

func DefaultVoted() uint64 {
	return 0
}
