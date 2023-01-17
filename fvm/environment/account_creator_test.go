package environment_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/fvm/environment"
	"github.com/koko1123/flow-go-1/fvm/state"
	"github.com/koko1123/flow-go-1/fvm/utils"
	"github.com/koko1123/flow-go-1/model/flow"
)

func Test_NewAccountCreator_NoError(t *testing.T) {
	view := utils.NewSimpleView()
	chain := flow.MonotonicEmulator.Chain()
	txnState := state.NewTransactionState(view, state.DefaultParameters())
	creator := environment.NewAddressGenerator(txnState, chain)
	require.NotNil(t, creator)
}

func Test_NewAccountCreator_GeneratingUpdatesState(t *testing.T) {
	view := utils.NewSimpleView()
	chain := flow.MonotonicEmulator.Chain()
	txnState := state.NewTransactionState(view, state.DefaultParameters())
	creator := environment.NewAddressGenerator(txnState, chain)
	_, err := creator.NextAddress()
	require.NoError(t, err)

	stateBytes, err := view.Get("", state.AddressStateKey)
	require.NoError(t, err)

	require.Equal(t, flow.BytesToAddress(stateBytes), flow.HexToAddress("01"))
}

func Test_NewAccountCreator_UsesLedgerState(t *testing.T) {
	view := utils.NewSimpleView()
	err := view.Set("", state.AddressStateKey, flow.HexToAddress("01").Bytes())
	require.NoError(t, err)

	chain := flow.MonotonicEmulator.Chain()
	txnState := state.NewTransactionState(view, state.DefaultParameters())
	creator := environment.NewAddressGenerator(txnState, chain)

	_, err = creator.NextAddress()
	require.NoError(t, err)

	stateBytes, err := view.Get("", state.AddressStateKey)
	require.NoError(t, err)

	require.Equal(t, flow.BytesToAddress(stateBytes), flow.HexToAddress("02"))
	// counts is one unit higher than returned index (index include zero, but counts starts from 1)
	require.Equal(t, uint64(2), creator.AddressCount())
}
