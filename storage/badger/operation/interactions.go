package operation

import (
	"github.com/koko1123/flow-go-1/engine/execution/state/delta"
	"github.com/koko1123/flow-go-1/model/flow"

	"github.com/dgraph-io/badger/v3"
)

func InsertExecutionStateInteractions(blockID flow.Identifier, interactions []*delta.Snapshot) func(*badger.Txn) error {
	return insert(makePrefix(codeExecutionStateInteractions, blockID), interactions)
}

func RetrieveExecutionStateInteractions(blockID flow.Identifier, interactions *[]*delta.Snapshot) func(*badger.Txn) error {
	return retrieve(makePrefix(codeExecutionStateInteractions, blockID), interactions)
}
