package consensus

import (
	"fmt"

	"github.com/gammazero/workerpool"
	"github.com/rs/zerolog"

	"github.com/koko1123/flow-go-1/consensus/hotstuff"
	"github.com/koko1123/flow-go-1/consensus/hotstuff/notifications/pubsub"
	"github.com/koko1123/flow-go-1/consensus/hotstuff/voteaggregator"
	"github.com/koko1123/flow-go-1/consensus/hotstuff/votecollector"
)

// NewVoteAggregator creates new VoteAggregator and recover the Forks' state with all pending block
func NewVoteAggregator(
	log zerolog.Logger,
	lowestRetainedView uint64,
	notifier hotstuff.Consumer,
	voteProcessorFactory hotstuff.VoteProcessorFactory,
	distributor *pubsub.FinalizationDistributor,
) (hotstuff.VoteAggregator, error) {

	createCollectorFactoryMethod := votecollector.NewStateMachineFactory(log, notifier, voteProcessorFactory.Create)
	voteCollectors := voteaggregator.NewVoteCollectors(log, lowestRetainedView, workerpool.New(4), createCollectorFactoryMethod)

	// initialize the vote aggregator
	aggregator, err := voteaggregator.NewVoteAggregator(log, notifier, lowestRetainedView, voteCollectors)
	if err != nil {
		return nil, fmt.Errorf("could not create vote aggregator: %w", err)
	}
	distributor.AddOnBlockFinalizedConsumer(aggregator.OnFinalizedBlock)

	return aggregator, nil
}
