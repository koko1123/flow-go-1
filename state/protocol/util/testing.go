package util

import (
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module"
	"github.com/koko1123/flow-go-1/module/metrics"
	modulemock "github.com/koko1123/flow-go-1/module/mock"
	"github.com/koko1123/flow-go-1/module/trace"
	"github.com/koko1123/flow-go-1/state/protocol"
	pbadger "github.com/koko1123/flow-go-1/state/protocol/badger"
	"github.com/koko1123/flow-go-1/state/protocol/events"
	mockprotocol "github.com/koko1123/flow-go-1/state/protocol/mock"
	"github.com/koko1123/flow-go-1/storage"
	"github.com/koko1123/flow-go-1/storage/util"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

// MockReceiptValidator returns a ReceiptValidator that accepts
// all receipts without performing any
// integrity checks.
func MockReceiptValidator() module.ReceiptValidator {
	validator := &modulemock.ReceiptValidator{}
	validator.On("Validate", mock.Anything).Return(nil)
	validator.On("ValidatePayload", mock.Anything).Return(nil)
	return validator
}

// MockBlockTimer returns BlockTimer that accepts all timestamps
// without performing any checks.
func MockBlockTimer() protocol.BlockTimer {
	blockTimer := &mockprotocol.BlockTimer{}
	blockTimer.On("Validate", mock.Anything, mock.Anything).Return(nil)
	return blockTimer
}

// MockSealValidator returns a SealValidator that accepts
// all seals without performing any
// integrity checks, returns first seal in block as valid one
func MockSealValidator(sealsDB storage.Seals) module.SealValidator {
	validator := &modulemock.SealValidator{}
	validator.On("Validate", mock.Anything).Return(
		func(candidate *flow.Block) *flow.Seal {
			if len(candidate.Payload.Seals) > 0 {
				return candidate.Payload.Seals[0]
			}
			last, _ := sealsDB.HighestInFork(candidate.Header.ParentID)
			return last
		},
		func(candidate *flow.Block) error {
			if len(candidate.Payload.Seals) > 0 {
				return nil
			}
			_, err := sealsDB.HighestInFork(candidate.Header.ParentID)
			return err
		}).Maybe()
	return validator
}

func RunWithBootstrapState(t testing.TB, rootSnapshot protocol.Snapshot, f func(*badger.DB, *pbadger.State)) {
	unittest.RunWithBadgerDB(t, func(db *badger.DB) {
		metrics := metrics.NewNoopCollector()
		headers, _, seals, _, _, blocks, setups, commits, statuses, results := util.StorageLayer(t, db)
		state, err := pbadger.Bootstrap(metrics, db, headers, seals, results, blocks, setups, commits, statuses, rootSnapshot)
		require.NoError(t, err)
		f(db, state)
	})
}

func RunWithFullProtocolState(t testing.TB, rootSnapshot protocol.Snapshot, f func(*badger.DB, *pbadger.MutableState)) {
	unittest.RunWithBadgerDB(t, func(db *badger.DB) {
		metrics := metrics.NewNoopCollector()
		tracer := trace.NewNoopTracer()
		consumer := events.NewNoop()
		headers, _, seals, index, payloads, blocks, setups, commits, statuses, results := util.StorageLayer(t, db)
		state, err := pbadger.Bootstrap(metrics, db, headers, seals, results, blocks, setups, commits, statuses, rootSnapshot)
		require.NoError(t, err)
		receiptValidator := MockReceiptValidator()
		sealValidator := MockSealValidator(seals)
		mockTimer := MockBlockTimer()
		fullState, err := pbadger.NewFullConsensusState(state, index, payloads, tracer, consumer, mockTimer, receiptValidator, sealValidator)
		require.NoError(t, err)
		f(db, fullState)
	})
}

func RunWithFullProtocolStateAndMetrics(t testing.TB, rootSnapshot protocol.Snapshot, metrics module.ComplianceMetrics, f func(*badger.DB, *pbadger.MutableState)) {
	unittest.RunWithBadgerDB(t, func(db *badger.DB) {
		tracer := trace.NewNoopTracer()
		consumer := events.NewNoop()
		headers, _, seals, index, payloads, blocks, setups, commits, statuses, results := util.StorageLayer(t, db)
		state, err := pbadger.Bootstrap(metrics, db, headers, seals, results, blocks, setups, commits, statuses, rootSnapshot)
		require.NoError(t, err)
		receiptValidator := MockReceiptValidator()
		sealValidator := MockSealValidator(seals)
		mockTimer := MockBlockTimer()
		fullState, err := pbadger.NewFullConsensusState(state, index, payloads, tracer, consumer, mockTimer, receiptValidator, sealValidator)
		require.NoError(t, err)
		f(db, fullState)
	})
}

func RunWithFullProtocolStateAndValidator(t testing.TB, rootSnapshot protocol.Snapshot, validator module.ReceiptValidator, f func(*badger.DB, *pbadger.MutableState)) {
	unittest.RunWithBadgerDB(t, func(db *badger.DB) {
		metrics := metrics.NewNoopCollector()
		tracer := trace.NewNoopTracer()
		consumer := events.NewNoop()
		headers, _, seals, index, payloads, blocks, setups, commits, statuses, results := util.StorageLayer(t, db)
		state, err := pbadger.Bootstrap(metrics, db, headers, seals, results, blocks, setups, commits, statuses, rootSnapshot)
		require.NoError(t, err)
		sealValidator := MockSealValidator(seals)
		mockTimer := MockBlockTimer()
		fullState, err := pbadger.NewFullConsensusState(state, index, payloads, tracer, consumer, mockTimer, validator, sealValidator)
		require.NoError(t, err)
		f(db, fullState)
	})
}

func RunWithFollowerProtocolState(t testing.TB, rootSnapshot protocol.Snapshot, f func(*badger.DB, *pbadger.FollowerState)) {
	unittest.RunWithBadgerDB(t, func(db *badger.DB) {
		metrics := metrics.NewNoopCollector()
		tracer := trace.NewNoopTracer()
		consumer := events.NewNoop()
		headers, _, seals, index, payloads, blocks, setups, commits, statuses, results := util.StorageLayer(t, db)
		state, err := pbadger.Bootstrap(metrics, db, headers, seals, results, blocks, setups, commits, statuses, rootSnapshot)
		require.NoError(t, err)
		mockTimer := MockBlockTimer()
		followerState, err := pbadger.NewFollowerState(state, index, payloads, tracer, consumer, mockTimer)
		require.NoError(t, err)
		f(db, followerState)
	})
}

func RunWithFullProtocolStateAndConsumer(t testing.TB, rootSnapshot protocol.Snapshot, consumer protocol.Consumer, f func(*badger.DB, *pbadger.MutableState)) {
	unittest.RunWithBadgerDB(t, func(db *badger.DB) {
		metrics := metrics.NewNoopCollector()
		tracer := trace.NewNoopTracer()
		headers, _, seals, index, payloads, blocks, setups, commits, statuses, results := util.StorageLayer(t, db)
		state, err := pbadger.Bootstrap(metrics, db, headers, seals, results, blocks, setups, commits, statuses, rootSnapshot)
		require.NoError(t, err)
		receiptValidator := MockReceiptValidator()
		sealValidator := MockSealValidator(seals)
		mockTimer := MockBlockTimer()
		fullState, err := pbadger.NewFullConsensusState(state, index, payloads, tracer, consumer, mockTimer, receiptValidator, sealValidator)
		require.NoError(t, err)
		f(db, fullState)
	})
}

func RunWithFollowerProtocolStateAndHeaders(t testing.TB, rootSnapshot protocol.Snapshot, f func(*badger.DB, *pbadger.FollowerState, storage.Headers, storage.Index)) {
	unittest.RunWithBadgerDB(t, func(db *badger.DB) {
		metrics := metrics.NewNoopCollector()
		tracer := trace.NewNoopTracer()
		consumer := events.NewNoop()
		headers, _, seals, index, payloads, blocks, setups, commits, statuses, results := util.StorageLayer(t, db)
		state, err := pbadger.Bootstrap(metrics, db, headers, seals, results, blocks, setups, commits, statuses, rootSnapshot)
		require.NoError(t, err)
		mockTimer := MockBlockTimer()
		followerState, err := pbadger.NewFollowerState(state, index, payloads, tracer, consumer, mockTimer)
		require.NoError(t, err)
		f(db, followerState, headers, index)
	})
}
