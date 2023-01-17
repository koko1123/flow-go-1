package bootstrap

import (
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/rs/zerolog"

	"github.com/koko1123/flow-go-1/engine/execution/state"
	"github.com/koko1123/flow-go-1/engine/execution/state/delta"
	"github.com/koko1123/flow-go-1/fvm"
	"github.com/koko1123/flow-go-1/ledger"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/storage"
	"github.com/koko1123/flow-go-1/storage/badger/operation"
)

// an increased limit for bootstrapping
const ledgerIntractionLimitNeededForBootstrapping = 1_000_000_000

type Bootstrapper struct {
	logger zerolog.Logger
}

func NewBootstrapper(logger zerolog.Logger) *Bootstrapper {
	return &Bootstrapper{
		logger: logger,
	}
}

// BootstrapLedger adds the above root account to the ledger and initializes execution node-only data
func (b *Bootstrapper) BootstrapLedger(
	ledger ledger.Ledger,
	servicePublicKey flow.AccountPublicKey,
	chain flow.Chain,
	opts ...fvm.BootstrapProcedureOption,
) (flow.StateCommitment, error) {
	view := delta.NewView(state.LedgerGetRegister(ledger, flow.StateCommitment(ledger.InitialState())))

	vm := fvm.NewVirtualMachine()

	ctx := fvm.NewContext(
		fvm.WithLogger(b.logger),
		fvm.WithMaxStateInteractionSize(ledgerIntractionLimitNeededForBootstrapping),
		fvm.WithChain(chain),
	)

	bootstrap := fvm.Bootstrap(
		servicePublicKey,
		opts...,
	)

	err := vm.Run(ctx, bootstrap, view)
	if err != nil {
		return flow.DummyStateCommitment, err
	}

	newStateCommitment, _, err := state.CommitDelta(ledger, view.Delta(), flow.StateCommitment(ledger.InitialState()))
	if err != nil {
		return flow.DummyStateCommitment, err
	}

	return newStateCommitment, nil
}

// IsBootstrapped returns whether the execution database has been bootstrapped, if yes, returns the
// root statecommitment
func (b *Bootstrapper) IsBootstrapped(db *badger.DB) (flow.StateCommitment, bool, error) {
	var commit flow.StateCommitment

	err := db.View(func(txn *badger.Txn) error {
		err := operation.LookupStateCommitment(flow.ZeroID, &commit)(txn)
		if err != nil {
			return fmt.Errorf("could not lookup state commitment: %w", err)
		}

		return nil
	})

	if errors.Is(err, storage.ErrNotFound) {
		return flow.DummyStateCommitment, false, nil
	}

	if err != nil {
		return flow.DummyStateCommitment, false, err
	}

	return commit, true, nil
}

func (b *Bootstrapper) BootstrapExecutionDatabase(db *badger.DB, commit flow.StateCommitment, genesis *flow.Header) error {

	err := operation.RetryOnConflict(db.Update, func(txn *badger.Txn) error {

		err := operation.InsertExecutedBlock(genesis.ID())(txn)
		if err != nil {
			return fmt.Errorf("could not index initial genesis execution block: %w", err)
		}

		err = operation.IndexStateCommitment(flow.ZeroID, commit)(txn)
		if err != nil {
			return fmt.Errorf("could not index void state commitment: %w", err)
		}

		err = operation.IndexStateCommitment(genesis.ID(), commit)(txn)
		if err != nil {
			return fmt.Errorf("could not index genesis state commitment: %w", err)
		}

		views := make([]*delta.Snapshot, 0)
		err = operation.InsertExecutionStateInteractions(genesis.ID(), views)(txn)
		if err != nil {
			return fmt.Errorf("could not bootstrap execution state interactions: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
