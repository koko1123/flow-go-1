// (c) 2019 Dapper Labs - ALL RIGHTS RESERVED

package storage

import (
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/storage/badger/transaction"
)

type EpochCommits interface {

	// StoreTx allows us to store a new epoch commit in a DB transaction while updating the cache.
	StoreTx(commit *flow.EpochCommit) func(*transaction.Tx) error

	// ByID will return the EpochCommit event by its ID.
	// Error returns:
	// * storage.ErrNotFound if no EpochCommit with the ID exists
	ByID(flow.Identifier) (*flow.EpochCommit, error)
}
