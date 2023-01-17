package chunkconsumer

import (
	"github.com/koko1123/flow-go-1/engine/verification/fetcher"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module"
)

// Worker receives job from job consumer and converts it back to Chunk
// for engine to process
type Worker struct {
	engine   fetcher.AssignedChunkProcessor
	consumer *ChunkConsumer
}

func NewWorker(engine fetcher.AssignedChunkProcessor) *Worker {
	return &Worker{
		engine: engine,
	}
}

// Run converts the job to Chunk, it's guaranteed to work, because
// ChunkJobs converted chunk into job symmetrically
func (w *Worker) Run(job module.Job) error {
	chunk, err := JobToChunkLocator(job)
	if err != nil {
		return err
	}
	w.engine.ProcessAssignedChunk(chunk)

	return nil
}

func (w *Worker) Notify(chunkLocatorID flow.Identifier) {
	jobID := locatorIDToJobID(chunkLocatorID)
	w.consumer.NotifyJobIsDone(jobID)
}
