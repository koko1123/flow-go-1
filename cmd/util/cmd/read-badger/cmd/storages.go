package cmd

import (
	"github.com/dgraph-io/badger/v3"

	"github.com/koko1123/flow-go-1/cmd/util/cmd/common"
	"github.com/koko1123/flow-go-1/storage"
)

func InitStorages() (*storage.All, *badger.DB) {
	db := common.InitStorage(flagDatadir)
	storages := common.InitStorages(db)
	return storages, db
}
