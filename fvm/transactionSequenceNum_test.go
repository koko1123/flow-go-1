package fvm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/fvm"
	"github.com/koko1123/flow-go-1/fvm/environment"
	"github.com/koko1123/flow-go-1/fvm/errors"
	"github.com/koko1123/flow-go-1/fvm/state"
	"github.com/koko1123/flow-go-1/fvm/tracing"
	"github.com/koko1123/flow-go-1/fvm/utils"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

func TestTransactionSequenceNumProcess(t *testing.T) {
	t.Run("sequence number update (happy path)", func(t *testing.T) {
		ledger := utils.NewSimpleView()
		txnState := state.NewTransactionState(ledger, state.DefaultParameters())
		accounts := environment.NewAccounts(txnState)

		// create an account
		address := flow.HexToAddress("1234")
		privKey, err := unittest.AccountKeyDefaultFixture()
		require.NoError(t, err)
		err = accounts.Create([]flow.AccountPublicKey{privKey.PublicKey(1000)}, address)
		require.NoError(t, err)

		tx := flow.TransactionBody{}
		tx.SetProposalKey(address, 0, 0)
		proc := fvm.Transaction(&tx, 0)

		seqChecker := fvm.TransactionSequenceNumberChecker{}
		err = seqChecker.CheckAndIncrementSequenceNumber(
			tracing.NewTracerSpan(),
			proc,
			txnState)
		require.NoError(t, err)

		// get fetch the sequence number and it should be updated
		key, err := accounts.GetPublicKey(address, 0)
		require.NoError(t, err)
		require.Equal(t, key.SeqNumber, uint64(1))
	})
	t.Run("invalid sequence number", func(t *testing.T) {
		ledger := utils.NewSimpleView()
		txnState := state.NewTransactionState(ledger, state.DefaultParameters())
		accounts := environment.NewAccounts(txnState)

		// create an account
		address := flow.HexToAddress("1234")
		privKey, err := unittest.AccountKeyDefaultFixture()
		require.NoError(t, err)
		err = accounts.Create([]flow.AccountPublicKey{privKey.PublicKey(1000)}, address)
		require.NoError(t, err)

		tx := flow.TransactionBody{}
		// invalid sequence number is 2
		tx.SetProposalKey(address, 0, 2)
		proc := fvm.Transaction(&tx, 0)

		seqChecker := fvm.TransactionSequenceNumberChecker{}
		err = seqChecker.CheckAndIncrementSequenceNumber(
			tracing.NewTracerSpan(),
			proc,
			txnState)
		require.Error(t, err)
		require.True(t, errors.HasErrorCode(err, errors.ErrCodeInvalidProposalSeqNumberError))

		// get fetch the sequence number and check it to be  unchanged
		key, err := accounts.GetPublicKey(address, 0)
		require.NoError(t, err)
		require.Equal(t, key.SeqNumber, uint64(0))
	})
	t.Run("invalid address", func(t *testing.T) {
		ledger := utils.NewSimpleView()
		txnState := state.NewTransactionState(ledger, state.DefaultParameters())
		accounts := environment.NewAccounts(txnState)

		// create an account
		address := flow.HexToAddress("1234")
		privKey, err := unittest.AccountKeyDefaultFixture()
		require.NoError(t, err)
		err = accounts.Create([]flow.AccountPublicKey{privKey.PublicKey(1000)}, address)
		require.NoError(t, err)

		tx := flow.TransactionBody{}
		// wrong address
		tx.SetProposalKey(flow.HexToAddress("2222"), 0, 0)
		proc := fvm.Transaction(&tx, 0)

		seqChecker := &fvm.TransactionSequenceNumberChecker{}
		err = seqChecker.CheckAndIncrementSequenceNumber(
			tracing.NewTracerSpan(),
			proc,
			txnState)
		require.Error(t, err)

		// get fetch the sequence number and check it to be unchanged
		key, err := accounts.GetPublicKey(address, 0)
		require.NoError(t, err)
		require.Equal(t, key.SeqNumber, uint64(0))
	})
}
