package ingestion

import (
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/model/messages"
	"github.com/koko1123/flow-go-1/module/mempool/stdmap"
)

type Deltas struct {
	*stdmap.Backend
}

// NewDeltas creates a new memory pool for state deltas
func NewDeltas(limit uint, opts ...stdmap.OptionFunc) (*Deltas, error) {
	s := &Deltas{
		Backend: stdmap.NewBackend(append(opts, stdmap.WithLimit(limit))...),
	}

	return s, nil
}

// Add adds an state deltas to the mempool.
func (s *Deltas) Add(delta *messages.ExecutionStateDelta) bool {
	// delta's ID is block's ID
	return s.Backend.Add(delta)
}

// Remove will remove a deltas by block ID.
func (s *Deltas) Remove(blockID flow.Identifier) bool {
	removed := s.Backend.Remove(blockID)
	return removed
}

// ByBlockID returns the state deltas for a block from the mempool.
func (s *Deltas) ByBlockID(blockID flow.Identifier) (*messages.ExecutionStateDelta, bool) {
	entity, exists := s.Backend.ByID(blockID)
	if !exists {
		return nil, false
	}
	delta := entity.(*messages.ExecutionStateDelta)
	return delta, true
}

// All returns all block Deltass from the pool.
func (s *Deltas) All() []*messages.ExecutionStateDelta {
	entities := s.Backend.All()
	deltas := make([]*messages.ExecutionStateDelta, 0, len(entities))
	for _, entity := range entities {
		deltas = append(deltas, entity.(*messages.ExecutionStateDelta))
	}
	return deltas
}
