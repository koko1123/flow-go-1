package unittest

import (
	"github.com/koko1123/flow-go-1/engine/execution"
	"github.com/koko1123/flow-go-1/engine/execution/state/delta"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module/mempool/entity"
	"github.com/koko1123/flow-go-1/utils/unittest"
	"github.com/onflow/flow-go/crypto"
)

func StateInteractionsFixture() *delta.SpockSnapshot {
	return delta.NewView(nil).Interactions()
}

func ComputationResultFixture(collectionsSignerIDs [][]flow.Identifier) *execution.ComputationResult {
	block := unittest.ExecutableBlockFixture(collectionsSignerIDs)
	startState := unittest.StateCommitmentFixture()
	block.StartState = &startState

	return ComputationResultForBlockFixture(block)
}

func ComputationResultForBlockFixture(
	completeBlock *entity.ExecutableBlock,
) *execution.ComputationResult {
	numChunks := len(completeBlock.CompleteCollections) + 1
	stateViews := make([]*delta.SpockSnapshot, numChunks)
	stateCommitments := make([]flow.StateCommitment, numChunks)
	proofs := make([][]byte, numChunks)
	events := make([]flow.EventsList, numChunks)
	eventHashes := make([]flow.Identifier, numChunks)
	spockHashes := make([]crypto.Signature, numChunks)
	for i := 0; i < numChunks; i++ {
		stateViews[i] = StateInteractionsFixture()
		stateCommitments[i] = *completeBlock.StartState
		proofs[i] = unittest.RandomBytes(6)
		events[i] = make(flow.EventsList, 0)
		eventHashes[i] = unittest.IdentifierFixture()
	}
	return &execution.ComputationResult{
		TransactionResultIndex: make([]int, numChunks),
		ExecutableBlock:        completeBlock,
		StateSnapshots:         stateViews,
		StateCommitments:       stateCommitments,
		Proofs:                 proofs,
		Events:                 events,
		EventsHashes:           eventHashes,
		SpockSignatures:        spockHashes,
	}
}
