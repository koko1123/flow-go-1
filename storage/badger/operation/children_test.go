package operation

import (
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

func TestBlockChildrenIndexUpdateLookup(t *testing.T) {
	unittest.RunWithBadgerDB(t, func(db *badger.DB) {
		blockID := unittest.IdentifierFixture()
		childrenIDs := unittest.IdentifierListFixture(8)
		var retrievedIDs flow.IdentifierList

		err := db.Update(InsertBlockChildren(blockID, childrenIDs))
		require.NoError(t, err)
		err = db.View(RetrieveBlockChildren(blockID, &retrievedIDs))
		require.NoError(t, err)
		assert.Equal(t, childrenIDs, retrievedIDs)

		altIDs := unittest.IdentifierListFixture(4)
		err = db.Update(UpdateBlockChildren(blockID, altIDs))
		require.NoError(t, err)
		err = db.View(RetrieveBlockChildren(blockID, &retrievedIDs))
		require.NoError(t, err)
		assert.Equal(t, altIDs, retrievedIDs)
	})
}
