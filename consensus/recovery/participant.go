package recovery

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/koko1123/flow-go-1/consensus/hotstuff"
	"github.com/koko1123/flow-go-1/consensus/hotstuff/model"
	"github.com/koko1123/flow-go-1/model/flow"
)

// Participant recovers the HotStuff state for a consensus participant.
// It reads the pending blocks from storage and pass them to the input Forks
// instance to recover its state from before the restart.
func Participant(
	log zerolog.Logger,
	forks hotstuff.Forks,
	voteAggregator hotstuff.VoteAggregator,
	validator hotstuff.Validator,
	finalized *flow.Header,
	pending []*flow.Header,
) error {
	return Recover(log, finalized, pending, validator, func(proposal *model.Proposal) error {
		// add it to forks
		err := forks.AddBlock(proposal.Block)
		if err != nil {
			return fmt.Errorf("could not add block to forks: %w", err)
		}

		// recovery the proposer's vote
		err = voteAggregator.AddBlock(proposal)
		if err != nil {
			return fmt.Errorf("could not process proposal %v: %w", proposal.Block.BlockID, err)
		}

		return nil
	})
}
