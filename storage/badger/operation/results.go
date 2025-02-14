package operation

import (
	"github.com/dgraph-io/badger/v3"

	"github.com/koko1123/flow-go-1/model/flow"
)

// InsertExecutionResult inserts an execution result by ID.
func InsertExecutionResult(result *flow.ExecutionResult) func(*badger.Txn) error {
	return insert(makePrefix(codeExecutionResult, result.ID()), result)
}

// BatchInsertExecutionResult inserts an execution result by ID.
func BatchInsertExecutionResult(result *flow.ExecutionResult) func(batch *badger.WriteBatch) error {
	return batchWrite(makePrefix(codeExecutionResult, result.ID()), result)
}

// RetrieveExecutionResult retrieves a transaction by fingerprint.
func RetrieveExecutionResult(resultID flow.Identifier, result *flow.ExecutionResult) func(*badger.Txn) error {
	return retrieve(makePrefix(codeExecutionResult, resultID), result)
}

// IndexExecutionResult inserts an execution result ID keyed by block ID
func IndexExecutionResult(blockID flow.Identifier, resultID flow.Identifier) func(*badger.Txn) error {
	return insert(makePrefix(codeIndexExecutionResultByBlock, blockID), resultID)
}

// ReindexExecutionResult updates mapping of an execution result ID keyed by block ID
func ReindexExecutionResult(blockID flow.Identifier, resultID flow.Identifier) func(*badger.Txn) error {
	return update(makePrefix(codeIndexExecutionResultByBlock, blockID), resultID)
}

// BatchIndexExecutionResult inserts an execution result ID keyed by block ID into a batch
func BatchIndexExecutionResult(blockID flow.Identifier, resultID flow.Identifier) func(batch *badger.WriteBatch) error {
	return batchWrite(makePrefix(codeIndexExecutionResultByBlock, blockID), resultID)
}

// LookupExecutionResult finds execution result ID by block
func LookupExecutionResult(blockID flow.Identifier, resultID *flow.Identifier) func(*badger.Txn) error {
	return retrieve(makePrefix(codeIndexExecutionResultByBlock, blockID), resultID)
}

// RemoveExecutionResultIndex removes execution result indexed by the given blockID
func RemoveExecutionResultIndex(blockID flow.Identifier) func(*badger.Txn) error {
	return remove(makePrefix(codeIndexExecutionResultByBlock, blockID))
}

// BatchRemoveExecutionResultIndex removes blockID-to-resultID index entries keyed by a blockID in a provided batch.
// No errors are expected during normal operation, even if no entries are matched.
// If Badger unexpectedly fails to process the request, the error is wrapped in a generic error and returned.
func BatchRemoveExecutionResultIndex(blockID flow.Identifier) func(*badger.WriteBatch) error {
	return batchRemove(makePrefix(codeIndexExecutionResultByBlock, blockID))
}
