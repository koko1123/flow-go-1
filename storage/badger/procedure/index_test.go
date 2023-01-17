package procedure

import (
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

func TestInsertRetrieveIndex(t *testing.T) {
	unittest.RunWithBadgerDB(t, func(db *badger.DB) {
		blockID := unittest.IdentifierFixture()
		index := unittest.IndexFixture()

		err := db.Update(InsertIndex(blockID, index))
		require.NoError(t, err)

		var retrieved flow.Index
		err = db.View(RetrieveIndex(blockID, &retrieved))
		require.NoError(t, err)

		require.Equal(t, index, &retrieved)
	})
}
