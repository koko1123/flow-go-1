package storage

import (
	"github.com/koko1123/flow-go-1/model/chunks"
)

type ChunksQueue interface {
	StoreChunkLocator(locator *chunks.Locator) (bool, error)

	LatestIndex() (uint64, error)

	AtIndex(index uint64) (*chunks.Locator, error)
}
