package common

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"

	"github.com/koko1123/flow-go-1/module/metrics"
	"github.com/koko1123/flow-go-1/state/protocol"
	protocolbadger "github.com/koko1123/flow-go-1/state/protocol/badger"
	"github.com/koko1123/flow-go-1/storage"
)

func InitProtocolState(db *badger.DB, storages *storage.All) (protocol.State, error) {
	metrics := &metrics.NoopCollector{}

	protocolState, err := protocolbadger.OpenState(
		metrics,
		db,
		storages.Headers,
		storages.Seals,
		storages.Results,
		storages.Blocks,
		storages.Setups,
		storages.EpochCommits,
		storages.Statuses,
	)

	if err != nil {
		return nil, fmt.Errorf("could not init protocol state: %w", err)
	}

	return protocolState, nil
}
