package storage

import (
	"github.com/koko1123/flow-go-1/model/flow"
)

// Guarantees represents persistent storage for collection guarantees.
type Guarantees interface {

	// Store inserts the collection guarantee.
	Store(guarantee *flow.CollectionGuarantee) error

	// ByCollectionID retrieves the collection guarantee by collection ID.
	ByCollectionID(collID flow.Identifier) (*flow.CollectionGuarantee, error)
}
