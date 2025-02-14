package common

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/rs/zerolog/log"

	"github.com/koko1123/flow-go-1/module/metrics"
	"github.com/koko1123/flow-go-1/storage"
	storagebadger "github.com/koko1123/flow-go-1/storage/badger"
	"github.com/koko1123/flow-go-1/storage/badger/operation"
)

func InitStorage(datadir string) *badger.DB {
	return InitStorageWithTruncate(datadir, false)
}

func InitStorageWithTruncate(datadir string, truncate bool) *badger.DB {
	opts := badger.
		DefaultOptions(datadir)

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal().Err(err).Msg("could not open key-value store")
	}

	// in order to void long iterations with big keys when initializing with an
	// already populated database, we bootstrap the initial maximum key size
	// upon starting
	err = operation.RetryOnConflict(db.Update, func(tx *badger.Txn) error {
		return operation.InitMax(tx)
	})
	if err != nil {
		log.Fatal().Err(err).Msg("could not initialize max tracker")
	}

	return db
}

func InitStorages(db *badger.DB) *storage.All {
	metrics := &metrics.NoopCollector{}

	return storagebadger.InitAll(metrics, db)
}
