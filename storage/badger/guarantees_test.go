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

func TestGuaranteeStoreRetrieve(t *testing.T) {
	unittest.RunWithBadgerDB(t, func(db *badger.DB) {
		metrics := metrics.NewNoopCollector()
		store := badgerstorage.NewGuarantees(metrics, db, 1000)

		// abiturary guarantees
		expected := unittest.CollectionGuaranteeFixture()

		// retrieve guarantee without stored
		_, err := store.ByCollectionID(expected.ID())
		require.True(t, errors.Is(err, storage.ErrNotFound))

		// store guarantee
		err = store.Store(expected)
		require.NoError(t, err)

		// retreive by coll idx
		actual, err := store.ByCollectionID(expected.ID())
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}
