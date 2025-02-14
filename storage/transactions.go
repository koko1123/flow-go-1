package storage

import (
	"github.com/koko1123/flow-go-1/model/flow"
)

// Transactions represents persistent storage for transactions.
type Transactions interface {

	// Store inserts the transaction, keyed by fingerprint. Duplicate transaction insertion is ignored
	Store(tx *flow.TransactionBody) error

	// ByID returns the transaction for the given fingerprint.
	ByID(txID flow.Identifier) (*flow.TransactionBody, error)
}
