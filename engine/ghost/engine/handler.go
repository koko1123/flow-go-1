package engine

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ghost "github.com/koko1123/flow-go-1/engine/ghost/protobuf"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/network"
	"github.com/koko1123/flow-go-1/network/channels"
)

// Handler handles the GRPC calls from a client
type Handler struct {
	log        zerolog.Logger
	conduitMap map[channels.Channel]network.Conduit
	msgChan    chan ghost.FlowMessage
	codec      network.Codec
}

var _ ghost.GhostNodeAPIServer = Handler{}

func NewHandler(log zerolog.Logger, conduitMap map[channels.Channel]network.Conduit, msgChan chan ghost.FlowMessage, codec network.Codec) *Handler {
	return &Handler{
		log:        log.With().Str("component", "ghost_engine_handler").Logger(),
		conduitMap: conduitMap,
		msgChan:    msgChan,
		codec:      codec,
	}
}

func (h Handler) SendEvent(_ context.Context, req *ghost.SendEventRequest) (*empty.Empty, error) {

	channelID := req.GetChannelId()

	// find the conduit for the channel
	conduit, found := h.conduitMap[channels.Channel(channelID)]

	if !found {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("conduit not found for given channel %v", channelID))
	}

	message := req.GetMessage()

	event, err := h.codec.Decode(message)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed to decode message")
	}

	targetIDs := req.GetTargetID()

	// Convert target ID to flow IDS
	var flowIDs = make([]flow.Identifier, len(targetIDs))
	for i, id := range targetIDs {
		flowIDs[i] = flow.HashToID(id)
	}

	h.log.Info().
		Interface("event", event).
		Str("flow_ids", fmt.Sprintf("%v", flowIDs)).
		Str("target_ids", fmt.Sprintf("%v", targetIDs)).
		Msg("sending message")

	// Submit the message over libp2p
	// TODO: there is an issue in the Publish method for the ghost node,
	// sometimes, it fails to deliver the message to the target without returning any error.
	// This becomes one of the big factors contributing to the tests flakeiness.
	err = conduit.Publish(event, flowIDs...)
	if err != nil {
		h.log.Error().
			Err(err).
			Interface("event", event).
			Str("flow_ids", fmt.Sprintf("%v", flowIDs)).
			Str("target_ids", fmt.Sprintf("%v", targetIDs)).
			Msg("error publishing message")

		return nil, err
	}

	return new(empty.Empty), nil
}

// Subscribe streams ALL the libp2p network messages over GRPC
func (h Handler) Subscribe(_ *ghost.SubscribeRequest, stream ghost.GhostNodeAPI_SubscribeServer) error {
	for {

		// read the network message from the channel
		flowMessage, ok := <-h.msgChan
		if !ok {
			return nil
		}

		// send it to the client
		err := stream.Send(&flowMessage)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to stream message: %v", err)
		}
	}
}
