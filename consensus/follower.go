package consensus

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/koko1123/flow-go-1/consensus/hotstuff"
	"github.com/koko1123/flow-go-1/consensus/hotstuff/follower"
	"github.com/koko1123/flow-go-1/consensus/hotstuff/validator"
	"github.com/koko1123/flow-go-1/consensus/recovery"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module"
	"github.com/koko1123/flow-go-1/storage"
)

// TODO: this needs to be integrated with proper configuration and bootstrapping.

func NewFollower(log zerolog.Logger, committee hotstuff.Committee, headers storage.Headers, updater module.Finalizer,
	verifier hotstuff.Verifier, notifier hotstuff.FinalizationConsumer, rootHeader *flow.Header,
	rootQC *flow.QuorumCertificate, finalized *flow.Header, pending []*flow.Header) (*hotstuff.FollowerLoop, error) {

	finalizer, err := newFinalizer(finalized, headers, updater, notifier, rootHeader, rootQC)
	if err != nil {
		return nil, fmt.Errorf("could not initialize finalizer: %w", err)
	}

	// initialize the Validator
	validator := validator.New(committee, finalizer, verifier)

	// recover the hotstuff state as a follower
	err = recovery.Follower(log, finalizer, validator, finalized, pending)
	if err != nil {
		return nil, fmt.Errorf("could not recover hotstuff follower state: %w", err)
	}

	// initialize the follower logic
	logic, err := follower.New(log, validator, finalizer)
	if err != nil {
		return nil, fmt.Errorf("could not create follower logic: %w", err)
	}

	// initialize the follower loop
	loop, err := hotstuff.NewFollowerLoop(log, logic)
	if err != nil {
		return nil, fmt.Errorf("could not create follower loop: %w", err)
	}

	return loop, nil
}
