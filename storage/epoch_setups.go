// (c) 2019 Dapper Labs - ALL RIGHTS RESERVED

package storage

import (
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/storage/badger/transaction"
)

type EpochSetups interface {

	// StoreTx allows us to store a new epoch setup in a DB transaction while going through the cache.
	StoreTx(*flow.EpochSetup) func(*transaction.Tx) error

	// ByID will return the EpochSetup event by its ID.
	// Error returns:
	// * storage.ErrNotFound if no EpochSetup with the ID exists
	ByID(flow.Identifier) (*flow.EpochSetup, error)
}
