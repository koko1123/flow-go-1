package badger_test

import (
	"errors"

	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/module/metrics"
	"github.com/koko1123/flow-go-1/storage"
	"github.com/koko1123/flow-go-1/utils/unittest"

	badgerstorage "github.com/koko1123/flow-go-1/storage/badger"
)

func TestPayloadStoreRetrieve(t *testing.T) {
	unittest.RunWithBadgerDB(t, func(db *badger.DB) {
		metrics := metrics.NewNoopCollector()

		index := badgerstorage.NewIndex(metrics, db)
		seals := badgerstorage.NewSeals(metrics, db)
		guarantees := badgerstorage.NewGuarantees(metrics, db, badgerstorage.DefaultCacheSize)
		results := badgerstorage.NewExecutionResults(metrics, db)
		receipts := badgerstorage.NewExecutionReceipts(metrics, db, results, badgerstorage.DefaultCacheSize)
		store := badgerstorage.NewPayloads(db, index, guarantees, seals, receipts, results)

		blockID := unittest.IdentifierFixture()
		expected := unittest.PayloadFixture(unittest.WithAllTheFixins)

		// store payload
		err := store.Store(blockID, &expected)
		require.NoError(t, err)

		// fetch payload
		payload, err := store.ByBlockID(blockID)
		require.NoError(t, err)
		require.Equal(t, &expected, payload)
	})
}

func TestPayloadRetreiveWithoutStore(t *testing.T) {
	unittest.RunWithBadgerDB(t, func(db *badger.DB) {
		metrics := metrics.NewNoopCollector()

		index := badgerstorage.NewIndex(metrics, db)
		seals := badgerstorage.NewSeals(metrics, db)
		guarantees := badgerstorage.NewGuarantees(metrics, db, badgerstorage.DefaultCacheSize)
		results := badgerstorage.NewExecutionResults(metrics, db)
		receipts := badgerstorage.NewExecutionReceipts(metrics, db, results, badgerstorage.DefaultCacheSize)
		store := badgerstorage.NewPayloads(db, index, guarantees, seals, receipts, results)

		blockID := unittest.IdentifierFixture()

		_, err := store.ByBlockID(blockID)
		require.True(t, errors.Is(err, storage.ErrNotFound))
	})
}
