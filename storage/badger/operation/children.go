package operation

import (
	"github.com/dgraph-io/badger/v3"

	"github.com/koko1123/flow-go-1/model/flow"
)

// InsertBlockChildren insert an index to lookup the direct child of a block by its ID
func InsertBlockChildren(blockID flow.Identifier, childrenIDs flow.IdentifierList) func(*badger.Txn) error {
	return insert(makePrefix(codeBlockChildren, blockID), childrenIDs)
}

// UpdateBlockChildren updates the children for a block.
func UpdateBlockChildren(blockID flow.Identifier, childrenIDs flow.IdentifierList) func(*badger.Txn) error {
	return update(makePrefix(codeBlockChildren, blockID), childrenIDs)
}

// RetrieveBlockChildren the child block ID by parent block ID
func RetrieveBlockChildren(blockID flow.Identifier, childrenIDs *flow.IdentifierList) func(*badger.Txn) error {
	return retrieve(makePrefix(codeBlockChildren, blockID), childrenIDs)
}
