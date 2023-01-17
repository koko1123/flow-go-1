package operation

import (
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

func TestTransactions(t *testing.T) {

	unittest.RunWithBadgerDB(t, func(db *badger.DB) {
		expected := unittest.TransactionFixture()
		err := db.Update(InsertTransaction(expected.ID(), &expected.TransactionBody))
		require.Nil(t, err)

		var actual flow.Transaction
		err = db.View(RetrieveTransaction(expected.ID(), &actual.TransactionBody))
		require.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}
