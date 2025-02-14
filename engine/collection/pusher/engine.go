// Package pusher implements an engine for providing access to resources held
// by the collection node, including collections, collection guarantees, and
// transactions.
package pusher

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/koko1123/flow-go-1/engine"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/model/flow/filter"
	"github.com/koko1123/flow-go-1/model/messages"
	"github.com/koko1123/flow-go-1/module"
	"github.com/koko1123/flow-go-1/module/metrics"
	"github.com/koko1123/flow-go-1/network"
	"github.com/koko1123/flow-go-1/network/channels"
	"github.com/koko1123/flow-go-1/state/protocol"
	"github.com/koko1123/flow-go-1/storage"
	"github.com/koko1123/flow-go-1/utils/logging"
)

// Engine is the collection pusher engine, which provides access to resources
// held by the collection node.
type Engine struct {
	unit         *engine.Unit
	log          zerolog.Logger
	engMetrics   module.EngineMetrics
	colMetrics   module.CollectionMetrics
	conduit      network.Conduit
	me           module.Local
	state        protocol.State
	collections  storage.Collections
	transactions storage.Transactions
}

func New(log zerolog.Logger, net network.Network, state protocol.State, engMetrics module.EngineMetrics, colMetrics module.CollectionMetrics, me module.Local, collections storage.Collections, transactions storage.Transactions) (*Engine, error) {
	e := &Engine{
		unit:         engine.NewUnit(),
		log:          log.With().Str("engine", "pusher").Logger(),
		engMetrics:   engMetrics,
		colMetrics:   colMetrics,
		me:           me,
		state:        state,
		collections:  collections,
		transactions: transactions,
	}

	conduit, err := net.Register(channels.PushGuarantees, e)
	if err != nil {
		return nil, fmt.Errorf("could not register for push protocol: %w", err)
	}
	e.conduit = conduit

	return e, nil
}

// Ready returns a ready channel that is closed once the engine has fully
// started.
func (e *Engine) Ready() <-chan struct{} {
	return e.unit.Ready()
}

// Done returns a done channel that is closed once the engine has fully stopped.
func (e *Engine) Done() <-chan struct{} {
	return e.unit.Done()
}

// SubmitLocal submits an event originating on the local node.
func (e *Engine) SubmitLocal(event interface{}) {
	e.unit.Launch(func() {
		err := e.process(e.me.NodeID(), event)
		if err != nil {
			engine.LogError(e.log, err)
		}
	})
}

// Submit submits the given event from the node with the given origin ID
// for processing in a non-blocking manner. It returns instantly and logs
// a potential processing error internally when done.
func (e *Engine) Submit(channel channels.Channel, originID flow.Identifier, event interface{}) {
	e.unit.Launch(func() {
		err := e.process(originID, event)
		if err != nil {
			engine.LogError(e.log, err)
		}
	})
}

// ProcessLocal processes an event originating on the local node.
func (e *Engine) ProcessLocal(event interface{}) error {
	return e.unit.Do(func() error {
		return e.process(e.me.NodeID(), event)
	})
}

// Process processes the given event from the node with the given origin ID in
// a blocking manner. It returns the potential processing error when done.
func (e *Engine) Process(channel channels.Channel, originID flow.Identifier, event interface{}) error {
	return e.unit.Do(func() error {
		return e.process(originID, event)
	})
}

// process processes events for the pusher engine on the collection node.
func (e *Engine) process(originID flow.Identifier, event interface{}) error {
	switch ev := event.(type) {
	case *messages.SubmitCollectionGuarantee:
		e.engMetrics.MessageReceived(metrics.EngineCollectionProvider, metrics.MessageSubmitGuarantee)
		defer e.engMetrics.MessageHandled(metrics.EngineCollectionProvider, metrics.MessageSubmitGuarantee)
		return e.onSubmitCollectionGuarantee(originID, ev)
	default:
		return fmt.Errorf("invalid event type (%T)", event)
	}
}

// onSubmitCollectionGuarantee handles submitting the given collection guarantee
// to consensus nodes.
func (e *Engine) onSubmitCollectionGuarantee(originID flow.Identifier, req *messages.SubmitCollectionGuarantee) error {
	if originID != e.me.NodeID() {
		return fmt.Errorf("invalid remote request to submit collection guarantee (from=%x)", originID)
	}

	return e.SubmitCollectionGuarantee(&req.Guarantee)
}

// SubmitCollectionGuarantee submits the collection guarantee to all consensus nodes.
func (e *Engine) SubmitCollectionGuarantee(guarantee *flow.CollectionGuarantee) error {
	consensusNodes, err := e.state.Final().Identities(filter.HasRole(flow.RoleConsensus))
	if err != nil {
		return fmt.Errorf("could not get consensus nodes: %w", err)
	}

	// NOTE: Consensus nodes do not broadcast guarantees among themselves, so it needs that
	// at least one collection node make a publish to all of them.
	err = e.conduit.Publish(guarantee, consensusNodes.NodeIDs()...)
	if err != nil {
		return fmt.Errorf("could not submit collection guarantee: %w", err)
	}

	e.engMetrics.MessageSent(metrics.EngineCollectionProvider, metrics.MessageCollectionGuarantee)

	e.log.Debug().
		Hex("guarantee_id", logging.ID(guarantee.ID())).
		Hex("ref_block_id", logging.ID(guarantee.ReferenceBlockID)).
		Msg("submitting collection guarantee")

	return nil
}
