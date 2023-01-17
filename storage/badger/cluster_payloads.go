package badger

import (
	"github.com/dgraph-io/badger/v3"

	"github.com/koko1123/flow-go-1/model/cluster"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module"
	"github.com/koko1123/flow-go-1/module/metrics"
	"github.com/koko1123/flow-go-1/storage/badger/operation"
	"github.com/koko1123/flow-go-1/storage/badger/procedure"
	"github.com/koko1123/flow-go-1/storage/badger/transaction"
)

// ClusterPayloads implements storage of block payloads for collection node
// cluster consensus.
type ClusterPayloads struct {
	db    *badger.DB
	cache *Cache
}

func NewClusterPayloads(cacheMetrics module.CacheMetrics, db *badger.DB) *ClusterPayloads {

	store := func(key interface{}, val interface{}) func(*transaction.Tx) error {
		blockID := key.(flow.Identifier)
		payload := val.(*cluster.Payload)
		return transaction.WithTx(procedure.InsertClusterPayload(blockID, payload))
	}

	retrieve := func(key interface{}) func(tx *badger.Txn) (interface{}, error) {
		blockID := key.(flow.Identifier)
		var payload cluster.Payload
		return func(tx *badger.Txn) (interface{}, error) {
			err := procedure.RetrieveClusterPayload(blockID, &payload)(tx)
			return &payload, err
		}
	}

	cp := &ClusterPayloads{
		db: db,
		cache: newCache(cacheMetrics, metrics.ResourceClusterPayload,
			withLimit(flow.DefaultTransactionExpiry*4),
			withStore(store),
			withRetrieve(retrieve)),
	}

	return cp
}

func (cp *ClusterPayloads) storeTx(blockID flow.Identifier, payload *cluster.Payload) func(*transaction.Tx) error {
	return cp.cache.PutTx(blockID, payload)
}
func (cp *ClusterPayloads) retrieveTx(blockID flow.Identifier) func(*badger.Txn) (*cluster.Payload, error) {
	return func(tx *badger.Txn) (*cluster.Payload, error) {
		val, err := cp.cache.Get(blockID)(tx)
		if err != nil {
			return nil, err
		}
		return val.(*cluster.Payload), nil
	}
}

func (cp *ClusterPayloads) Store(blockID flow.Identifier, payload *cluster.Payload) error {
	return operation.RetryOnConflictTx(cp.db, transaction.Update, cp.storeTx(blockID, payload))
}

func (cp *ClusterPayloads) ByBlockID(blockID flow.Identifier) (*cluster.Payload, error) {
	tx := cp.db.NewTransaction(false)
	defer tx.Discard()
	return cp.retrieveTx(blockID)(tx)
}
