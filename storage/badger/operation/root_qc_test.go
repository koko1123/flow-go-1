package operation

import (
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

func TestInsertRetrieveRootQC(t *testing.T) {
	qc := unittest.QuorumCertificateFixture()

	unittest.RunWithBadgerDB(t, func(db *badger.DB) {
		err := db.Update(InsertRootQuorumCertificate(qc))
		require.NoError(t, err)

		// should be able to retrieve
		var retrieved flow.QuorumCertificate
		err = db.View(RetrieveRootQuorumCertificate(&retrieved))
		require.NoError(t, err)
		assert.Equal(t, qc, &retrieved)

		// should not be able to overwrite
		qc2 := unittest.QuorumCertificateFixture()
		err = db.Update(InsertRootQuorumCertificate(qc2))
		require.Error(t, err)
	})
}
