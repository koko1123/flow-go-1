package environment

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/fvm/state"
	"github.com/koko1123/flow-go-1/fvm/tracing"
	"github.com/koko1123/flow-go-1/fvm/utils"
)

func TestUUIDs_GetAndSetUUID(t *testing.T) {
	view := utils.NewSimpleView()
	txnState := state.NewTransactionState(view, state.DefaultParameters())
	uuidsA := NewUUIDGenerator(
		tracing.NewTracerSpan(),
		NewMeter(txnState),
		txnState)

	uuid, err := uuidsA.getUUID() // start from zero
	require.NoError(t, err)
	require.Equal(t, uint64(0), uuid)

	err = uuidsA.setUUID(5)
	require.NoError(t, err)

	// create new UUIDs instance
	uuidsB := NewUUIDGenerator(
		tracing.NewTracerSpan(),
		NewMeter(txnState),
		txnState)
	uuid, err = uuidsB.getUUID() // should read saved value
	require.NoError(t, err)

	require.Equal(t, uint64(5), uuid)
}

func Test_GenerateUUID(t *testing.T) {
	view := utils.NewSimpleView()
	txnState := state.NewTransactionState(view, state.DefaultParameters())
	genA := NewUUIDGenerator(
		tracing.NewTracerSpan(),
		NewMeter(txnState),
		txnState)

	uuidA, err := genA.GenerateUUID()
	require.NoError(t, err)
	uuidB, err := genA.GenerateUUID()
	require.NoError(t, err)
	uuidC, err := genA.GenerateUUID()
	require.NoError(t, err)

	require.Equal(t, uint64(0), uuidA)
	require.Equal(t, uint64(1), uuidB)
	require.Equal(t, uint64(2), uuidC)

	// Create new generator instance from same ledger
	genB := NewUUIDGenerator(
		tracing.NewTracerSpan(),
		NewMeter(txnState),
		txnState)

	uuidD, err := genB.GenerateUUID()
	require.NoError(t, err)
	uuidE, err := genB.GenerateUUID()
	require.NoError(t, err)
	uuidF, err := genB.GenerateUUID()
	require.NoError(t, err)

	require.Equal(t, uint64(3), uuidD)
	require.Equal(t, uint64(4), uuidE)
	require.Equal(t, uint64(5), uuidF)
}
