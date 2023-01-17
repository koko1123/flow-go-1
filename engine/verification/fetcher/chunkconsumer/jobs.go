package chunkconsumer

import (
	"fmt"

	"github.com/koko1123/flow-go-1/module"
	"github.com/koko1123/flow-go-1/storage"
)

// ChunkJobs wraps the storage layer to provide an abstraction for consumers to read jobs.
type ChunkJobs struct {
	locators storage.ChunksQueue
}

func (j *ChunkJobs) AtIndex(index uint64) (module.Job, error) {
	locator, err := j.locators.AtIndex(index)
	if err != nil {
		return nil, fmt.Errorf("could not read chunk: %w", err)
	}
	return ChunkLocatorToJob(locator), nil
}

func (j *ChunkJobs) Head() (uint64, error) {
	return j.locators.LatestIndex()
}
