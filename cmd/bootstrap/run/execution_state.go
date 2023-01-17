package run

import (
	"math"

	"github.com/rs/zerolog"
	"go.uber.org/atomic"

	"github.com/koko1123/flow-go-1/crypto"
	"github.com/koko1123/flow-go-1/crypto/hash"
	"github.com/koko1123/flow-go-1/engine/execution/state/bootstrap"
	"github.com/koko1123/flow-go-1/fvm"
	"github.com/koko1123/flow-go-1/ledger/common/pathfinder"
	"github.com/koko1123/flow-go-1/ledger/complete"
	ledger "github.com/koko1123/flow-go-1/ledger/complete"
	"github.com/koko1123/flow-go-1/ledger/complete/wal"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module/metrics"
)

// NOTE: this is now unused and should become part of another tool.
func GenerateServiceAccountPrivateKey(seed []byte) (flow.AccountPrivateKey, error) {
	priv, err := crypto.GeneratePrivateKey(crypto.ECDSASecp256k1, seed)
	if err != nil {
		return flow.AccountPrivateKey{}, err
	}

	return flow.AccountPrivateKey{
		PrivateKey: priv,
		SignAlgo:   crypto.ECDSASecp256k1,
		HashAlgo:   hash.SHA2_256,
	}, nil
}

// NOTE: this is now unused and should become part of another tool.
func GenerateExecutionState(
	dbDir string,
	accountKey flow.AccountPublicKey,
	chain flow.Chain,
	bootstrapOptions ...fvm.BootstrapProcedureOption,
) (flow.StateCommitment, error) {
	const (
		capacity           = 100
		checkpointDistance = math.MaxInt // A large number to prevent checkpoint creation.
		checkpointsToKeep  = 1
	)

	metricsCollector := &metrics.NoopCollector{}

	diskWal, err := wal.NewDiskWAL(zerolog.Nop(), nil, metricsCollector, dbDir, capacity, pathfinder.PathByteSize, wal.SegmentSize)
	if err != nil {
		return flow.DummyStateCommitment, err
	}

	ledgerStorage, err := ledger.NewLedger(diskWal, capacity, metricsCollector, zerolog.Nop(), ledger.DefaultPathFinderVersion)
	if err != nil {
		return flow.DummyStateCommitment, err
	}

	compactor, err := complete.NewCompactor(ledgerStorage, diskWal, zerolog.Nop(), capacity, checkpointDistance, checkpointsToKeep, atomic.NewBool(false))
	if err != nil {
		return flow.DummyStateCommitment, err
	}
	<-compactor.Ready()

	defer func() {
		<-ledgerStorage.Done()
		<-compactor.Done()
	}()

	return bootstrap.NewBootstrapper(
		zerolog.Nop()).BootstrapLedger(
		ledgerStorage,
		accountKey,
		chain,
		bootstrapOptions...,
	)
}
