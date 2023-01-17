package run

import (
	"fmt"

	"github.com/koko1123/flow-go-1/model/flow"
)

func GenerateRootSeal(result *flow.ExecutionResult) (*flow.Seal, error) {
	finalState, err := result.FinalStateCommitment()
	if err != nil {
		return nil, fmt.Errorf("generating root seal failed: %w", err)
	}
	seal := &flow.Seal{
		BlockID:    result.BlockID,
		ResultID:   result.ID(),
		FinalState: finalState,
	}
	return seal, nil
}
