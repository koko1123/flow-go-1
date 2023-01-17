package committer_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/engine/execution/computation/committer"
	fvmUtils "github.com/koko1123/flow-go-1/fvm/utils"
	led "github.com/koko1123/flow-go-1/ledger"
	ledgermock "github.com/koko1123/flow-go-1/ledger/mock"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module/trace"
	utils "github.com/koko1123/flow-go-1/utils/unittest"
)

func TestLedgerViewCommitter(t *testing.T) {

	t.Run("calls to set and prove", func(t *testing.T) {

		ledger := new(ledgermock.Ledger)
		com := committer.NewLedgerViewCommitter(ledger, trace.NewNoopTracer())

		var expectedStateCommitment led.State
		copy(expectedStateCommitment[:], []byte{1, 2, 3})
		ledger.On("Set", mock.Anything).
			Return(expectedStateCommitment, nil, nil).
			Once()

		expectedProof := led.Proof([]byte{2, 3, 4})
		ledger.On("Prove", mock.Anything).
			Return(expectedProof, nil).
			Once()

		view := fvmUtils.NewSimpleView()

		err := view.Set(
			"owner",
			"key",
			[]byte{1},
		)
		require.NoError(t, err)

		newState, proof, _, err := com.CommitView(view, utils.StateCommitmentFixture())
		require.NoError(t, err)
		require.Equal(t, flow.StateCommitment(expectedStateCommitment), newState)
		require.Equal(t, []uint8(expectedProof), proof)
	})

}
