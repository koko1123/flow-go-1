package consensus

import (
	"fmt"

	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module"
	"github.com/koko1123/flow-go-1/module/mempool"
	"github.com/koko1123/flow-go-1/module/metrics"
	"github.com/koko1123/flow-go-1/storage"
)

// CleanupFunc is called after a block was finalized to allow other components
// to execute cleanup operations.
type CleanupFunc func(blockID flow.Identifier) error

func CleanupNothing() CleanupFunc {
	return func(flow.Identifier) error {
		return nil
	}
}

func CleanupMempools(
	collector module.MempoolMetrics,
	spans module.ConsensusMetrics,
	payloads storage.Payloads,
	guarantees mempool.Guarantees,
	seals mempool.IncorporatedResultSeals) CleanupFunc {

	return func(blockID flow.Identifier) error {

		payload, err := payloads.ByBlockID(blockID)
		if err != nil {
			return fmt.Errorf("could not retrieve  payload (%x): %w", blockID, err)
		}

		for _, guarantee := range payload.Guarantees {
			_ = guarantees.Remove(guarantee.ID())
		}

		collector.MempoolEntries(metrics.ResourceGuarantee, guarantees.Size())

		for _, seal := range payload.Seals {
			_ = seals.Remove(seal.ID())
		}

		collector.MempoolEntries(metrics.ResourceSeal, seals.Size())

		return nil
	}
}
