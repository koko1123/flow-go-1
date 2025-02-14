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

// TestEpochSetupStoreAndRetrieve tests that a setup can be stored, retrieved and attempted to be stored again without an error
func TestEpochSetupStoreAndRetrieve(t *testing.T) {
	unittest.RunWithBadgerDB(t, func(db *badger.DB) {
		metrics := metrics.NewNoopCollector()
		store := badgerstorage.NewEpochSetups(metrics, db)

		// attempt to get a setup that doesn't exist
		_, err := store.ByID(unittest.IdentifierFixture())
		assert.True(t, errors.Is(err, storage.ErrNotFound))

		// store a setup in db
		expected := unittest.EpochSetupFixture()
		err = operation.RetryOnConflictTx(db, transaction.Update, store.StoreTx(expected))
		require.NoError(t, err)

		// retrieve the setup by ID
		actual, err := store.ByID(expected.ID())
		require.NoError(t, err)
		assert.Equal(t, expected, actual)

		// test storing same epoch setup
		err = operation.RetryOnConflictTx(db, transaction.Update, store.StoreTx(expected))
		require.NoError(t, err)
	})
}
