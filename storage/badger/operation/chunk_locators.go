package operation

import (
	"github.com/dgraph-io/badger/v3"

	"github.com/koko1123/flow-go-1/model/chunks"
	"github.com/koko1123/flow-go-1/model/flow"
)

func InsertChunkLocator(locator *chunks.Locator) func(*badger.Txn) error {
	return insert(makePrefix(codeChunk, locator.ID()), locator)
}

func RetrieveChunkLocator(locatorID flow.Identifier, locator *chunks.Locator) func(*badger.Txn) error {
	return retrieve(makePrefix(codeChunk, locatorID), locator)
}
