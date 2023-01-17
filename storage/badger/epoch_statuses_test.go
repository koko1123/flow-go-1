package badger_test

import (
	"errors"
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/module/metrics"
	"github.com/koko1123/flow-go-1/storage"
	"github.com/koko1123/flow-go-1/utils/unittest"

	badgerstorage "github.com/koko1123/flow-go-1/storage/badger"
	"github.com/koko1123/flow-go-1/storage/badger/operation"
	"github.com/koko1123/flow-go-1/storage/badger/transaction"
)

func TestEpochStatusesStoreAndRetrieve(t *testing.T) {
	unittest.RunWithBadgerDB(t, func(db *badger.DB) {
		metrics := metrics.NewNoopCollector()
		store := badgerstorage.NewEpochStatuses(metrics, db)

		blockID := unittest.IdentifierFixture()
		expected := unittest.EpochStatusFixture()

		_, err := store.ByBlockID(unittest.IdentifierFixture())
		assert.True(t, errors.Is(err, storage.ErrNotFound))

		// store epoch status
		err = operation.RetryOnConflictTx(db, transaction.Update, store.StoreTx(blockID, expected))
		require.NoError(t, err)

		// retreive status
		actual, err := store.ByBlockID(blockID)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}
