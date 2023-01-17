package synchronization

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/koko1123/flow-go-1/engine"
	commonsync "github.com/koko1123/flow-go-1/engine/common/synchronization"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/model/messages"
	"github.com/koko1123/flow-go-1/module"
	"github.com/koko1123/flow-go-1/module/chainsync"
	"github.com/koko1123/flow-go-1/module/lifecycle"
	"github.com/koko1123/flow-go-1/module/metrics"
	"github.com/koko1123/flow-go-1/network"
	"github.com/koko1123/flow-go-1/network/channels"
	"github.com/koko1123/flow-go-1/state/cluster"
	"github.com/koko1123/flow-go-1/storage"
)

// defaultSyncRequestQueueCapacity maximum capacity of sync requests queue
const defaultSyncRequestQueueCapacity = 500

// defaultSyncRequestQueueCapacity maximum capacity of range requests queue
const defaultRangeRequestQueueCapacity = 500

// defaultSyncRequestQueueCapacity maximum capacity of batch requests queue
const defaultBatchRequestQueueCapacity = 500

// defaultEngineRequestsWorkers number of workers to dispatch events for requests
const defaultEngineRequestsWorkers = 8

type RequestHandlerEngine struct {
	unit *engine.Unit
	lm   *lifecycle.LifecycleManager

	me      module.Local
	log     zerolog.Logger
	metrics module.EngineMetrics

	blocks storage.ClusterBlocks
	core   module.SyncCore
	state  cluster.State
	con    network.Conduit // used for sending responses to requesters

	pendingSyncRequests   engine.MessageStore    // message store for *message.SyncRequest
	pendingBatchRequests  engine.MessageStore    // message store for *message.BatchRequest
	pendingRangeRequests  engine.MessageStore    // message store for *message.RangeRequest
	requestMessageHandler *engine.MessageHandler // message handler responsible for request processing
}

func NewRequestHandlerEngine(
	log zerolog.Logger,
	metrics module.EngineMetrics,
	con network.Conduit,
	me module.Local,
	blocks storage.ClusterBlocks,
	core module.SyncCore,
	state cluster.State,
) *RequestHandlerEngine {
	r := &RequestHandlerEngine{
		unit:    engine.NewUnit(),
		lm:      lifecycle.NewLifecycleManager(),
		me:      me,
		log:     log.With().Str("engine", "cluster_synchronization").Logger(),
		metrics: metrics,
		blocks:  blocks,
		core:    core,
		state:   state,
		con:     con,
	}

	r.setupRequestMessageHandler()

	return r
}

// SubmitLocal submits an event originating on the local node.
func (r *RequestHandlerEngine) SubmitLocal(event interface{}) {
	err := r.ProcessLocal(event)
	if err != nil {
		r.log.Fatal().Err(err).Msg("internal error processing event")
	}
}

// Submit submits the given event from the node with the given origin ID
// for processing in a non-blocking manner. It returns instantly and logs
// a potential processing error internally when done.
func (r *RequestHandlerEngine) Submit(channel channels.Channel, originID flow.Identifier, event interface{}) {
	err := r.Process(channel, originID, event)
	if err != nil {
		r.log.Fatal().Err(err).Msg("internal error processing event")
	}
}

// ProcessLocal processes an event originating on the local node.
func (r *RequestHandlerEngine) ProcessLocal(event interface{}) error {
	return r.process(r.me.NodeID(), event)
}

// Process processes the given event from the node with the given origin ID in
// a blocking manner. It returns the potential processing error when done.
func (r *RequestHandlerEngine) Process(channel channels.Channel, originID flow.Identifier, event interface{}) error {
	err := r.process(originID, event)
	if err != nil {
		if engine.IsIncompatibleInputTypeError(err) {
			r.log.Warn().Msgf("%v delivered unsupported message %T through %v", originID, event, channel)
			return nil
		}
		return fmt.Errorf("unexpected error while processing engine message: %w", err)
	}
	return nil
}

// process processes events for the synchronization request handler engine.
// Error returns:
//   - IncompatibleInputTypeError if input has unexpected type
//   - All other errors are potential symptoms of internal state corruption or bugs (fatal).
func (r *RequestHandlerEngine) process(originID flow.Identifier, event interface{}) error {
	return r.requestMessageHandler.Process(originID, event)
}

// setupRequestMessageHandler initializes the inbound queues and the MessageHandler for UNTRUSTED requests.
func (r *RequestHandlerEngine) setupRequestMessageHandler() {
	// RequestHeap deduplicates requests by keeping only one sync request for each requester.
	r.pendingSyncRequests = commonsync.NewRequestHeap(defaultSyncRequestQueueCapacity)
	r.pendingRangeRequests = commonsync.NewRequestHeap(defaultRangeRequestQueueCapacity)
	r.pendingBatchRequests = commonsync.NewRequestHeap(defaultBatchRequestQueueCapacity)

	// define message queueing behaviour
	r.requestMessageHandler = engine.NewMessageHandler(
		r.log,
		engine.NewNotifier(),
		engine.Pattern{
			Match: func(msg *engine.Message) bool {
				_, ok := msg.Payload.(*messages.SyncRequest)
				if ok {
					r.metrics.MessageReceived(metrics.EngineClusterSynchronization, metrics.MessageSyncRequest)
				}
				return ok
			},
			Store: r.pendingSyncRequests,
		},
		engine.Pattern{
			Match: func(msg *engine.Message) bool {
				_, ok := msg.Payload.(*messages.RangeRequest)
				if ok {
					r.metrics.MessageReceived(metrics.EngineClusterSynchronization, metrics.MessageRangeRequest)
				}
				return ok
			},
			Store: r.pendingRangeRequests,
		},
		engine.Pattern{
			Match: func(msg *engine.Message) bool {
				_, ok := msg.Payload.(*messages.BatchRequest)
				if ok {
					r.metrics.MessageReceived(metrics.EngineClusterSynchronization, metrics.MessageBatchRequest)
				}
				return ok
			},
			Store: r.pendingBatchRequests,
		},
	)
}

// onSyncRequest processes an outgoing handshake; if we have a higher height, we
// inform the other node of it, so they can organize their block downloads. If
// we have a lower height, we add the difference to our own download queue.
func (r *RequestHandlerEngine) onSyncRequest(originID flow.Identifier, req *messages.SyncRequest) error {
	final, err := r.state.Final().Head()
	if err != nil {
		return fmt.Errorf("could not get last finalized header: %w", err)
	}

	// queue any missing heights as needed
	r.core.HandleHeight(final, req.Height)

	// don't bother sending a response if we're within tolerance or if we're
	// behind the requester
	if r.core.WithinTolerance(final, req.Height) || req.Height > final.Height {
		return nil
	}

	// if we're sufficiently ahead of the requester, send a response
	res := &messages.SyncResponse{
		Height: final.Height,
		Nonce:  req.Nonce,
	}
	err = r.con.Unicast(res, originID)
	if err != nil {
		r.log.Warn().Err(err).Msg("sending sync response failed")
		return nil
	}
	r.metrics.MessageSent(metrics.EngineClusterSynchronization, metrics.MessageSyncResponse)

	return nil
}

// onRangeRequest processes a request for a range of blocks by height.
func (r *RequestHandlerEngine) onRangeRequest(originID flow.Identifier, req *messages.RangeRequest) error {
	r.log.Debug().Str("origin_id", originID.String()).Msg("received new range request")
	// get the latest final state to know if we can fulfill the request
	head, err := r.state.Final().Head()
	if err != nil {
		return fmt.Errorf("could not get last finalized header: %w", err)
	}

	// if we don't have anything to send, we can bail right away
	if head.Height < req.FromHeight || req.FromHeight > req.ToHeight {
		return nil
	}

	// enforce client-side max request size
	var maxSize uint
	// TODO: clean up this logic
	if core, ok := r.core.(*chainsync.Core); ok {
		maxSize = core.Config.MaxSize
	} else {
		maxSize = chainsync.DefaultConfig().MaxSize
	}
	maxHeight := req.FromHeight + uint64(maxSize)
	if maxHeight < req.ToHeight {
		r.log.Warn().
			Uint64("from", req.FromHeight).
			Uint64("to", req.ToHeight).
			Uint64("size", (req.ToHeight-req.FromHeight)+1).
			Uint("max_size", maxSize).
			Msg("range request is too large")

		req.ToHeight = maxHeight
	}

	// get all of the blocks, one by one
	blocks := make([]messages.UntrustedClusterBlock, 0, req.ToHeight-req.FromHeight+1)
	for height := req.FromHeight; height <= req.ToHeight; height++ {
		block, err := r.blocks.ByHeight(height)
		if errors.Is(err, storage.ErrNotFound) {
			r.log.Error().Uint64("height", height).Msg("skipping unknown heights")
			break
		}
		if err != nil {
			return fmt.Errorf("could not get block for height (%d): %w", height, err)
		}
		blocks = append(blocks, messages.UntrustedClusterBlockFromInternal(block))
	}

	// if there are no blocks to send, skip network message
	if len(blocks) == 0 {
		r.log.Debug().Msg("skipping empty range response")
		return nil
	}

	// send the response
	res := &messages.ClusterBlockResponse{
		Nonce:  req.Nonce,
		Blocks: blocks,
	}
	err = r.con.Unicast(res, originID)
	if err != nil {
		r.log.Warn().Err(err).Hex("origin_id", originID[:]).Msg("sending range response failed")
		return nil
	}
	r.metrics.MessageSent(metrics.EngineClusterSynchronization, metrics.MessageBlockResponse)

	return nil
}

// onBatchRequest processes a request for a specific block by block ID.
func (r *RequestHandlerEngine) onBatchRequest(originID flow.Identifier, req *messages.BatchRequest) error {
	r.log.Debug().Str("origin_id", originID.String()).Msg("received new batch request")
	// we should bail and send nothing on empty request
	if len(req.BlockIDs) == 0 {
		return nil
	}

	// TODO: clean up this logic
	var maxSize uint
	if core, ok := r.core.(*chainsync.Core); ok {
		maxSize = core.Config.MaxSize
	} else {
		maxSize = chainsync.DefaultConfig().MaxSize
	}

	if len(req.BlockIDs) > int(maxSize) {
		r.log.Warn().
			Int("size", len(req.BlockIDs)).
			Uint("max_size", maxSize).
			Msg("batch request is too large")
	}

	// deduplicate the block IDs in the batch request
	blockIDs := make(map[flow.Identifier]struct{})
	for _, blockID := range req.BlockIDs {
		blockIDs[blockID] = struct{}{}

		// enforce client-side max request size
		if len(blockIDs) == int(maxSize) {
			break
		}
	}

	// try to get all the blocks by ID
	blocks := make([]messages.UntrustedClusterBlock, 0, len(blockIDs))
	for blockID := range blockIDs {
		block, err := r.blocks.ByID(blockID)
		if errors.Is(err, storage.ErrNotFound) {
			r.log.Debug().Hex("block_id", blockID[:]).Msg("skipping unknown block")
			continue
		}
		if err != nil {
			return fmt.Errorf("could not get block by ID (%s): %w", blockID, err)
		}
		blocks = append(blocks, messages.UntrustedClusterBlockFromInternal(block))
	}

	// if there are no blocks to send, skip network message
	if len(blocks) == 0 {
		r.log.Debug().Msg("skipping empty batch response")
		return nil
	}

	// send the response
	res := &messages.ClusterBlockResponse{
		Nonce:  req.Nonce,
		Blocks: blocks,
	}
	err := r.con.Unicast(res, originID)
	if err != nil {
		r.log.Warn().Err(err).Hex("origin_id", originID[:]).Msg("sending batch response failed")
		return nil
	}
	r.metrics.MessageSent(metrics.EngineClusterSynchronization, metrics.MessageBlockResponse)

	return nil
}

// processAvailableRequests is processor of pending events which drives events from networking layer to business logic.
func (r *RequestHandlerEngine) processAvailableRequests() error {
	for {
		select {
		case <-r.unit.Quit():
			return nil
		default:
		}

		msg, ok := r.pendingSyncRequests.Get()
		if ok {
			err := r.onSyncRequest(msg.OriginID, msg.Payload.(*messages.SyncRequest))
			if err != nil {
				return fmt.Errorf("processing sync request failed: %w", err)
			}
			continue
		}

		msg, ok = r.pendingRangeRequests.Get()
		if ok {
			err := r.onRangeRequest(msg.OriginID, msg.Payload.(*messages.RangeRequest))
			if err != nil {
				return fmt.Errorf("processing range request failed: %w", err)
			}
			continue
		}

		msg, ok = r.pendingBatchRequests.Get()
		if ok {
			err := r.onBatchRequest(msg.OriginID, msg.Payload.(*messages.BatchRequest))
			if err != nil {
				return fmt.Errorf("processing batch request failed: %w", err)
			}
			continue
		}

		// when there is no more messages in the queue, back to the loop to wait
		// for the next incoming message to arrive.
		return nil
	}
}

// requestProcessingLoop is a separate goroutine that performs processing of queued requests
func (r *RequestHandlerEngine) requestProcessingLoop() {
	notifier := r.requestMessageHandler.GetNotifier()
	for {
		select {
		case <-r.unit.Quit():
			return
		case <-notifier:
			err := r.processAvailableRequests()
			if err != nil {
				r.log.Fatal().Err(err).Msg("internal error processing queued requests")
			}
		}
	}
}

// Ready returns a ready channel that is closed once the engine has fully started.
func (r *RequestHandlerEngine) Ready() <-chan struct{} {
	r.lm.OnStart(func() {
		for i := 0; i < defaultEngineRequestsWorkers; i++ {
			r.unit.Launch(r.requestProcessingLoop)
		}
	})
	return r.lm.Started()
}

// Done returns a done channel that is closed once the engine has fully stopped.
func (r *RequestHandlerEngine) Done() <-chan struct{} {
	r.lm.OnStop(func() {
		// wait for all request processing workers to exit
		<-r.unit.Done()
	})
	return r.lm.Stopped()
}
