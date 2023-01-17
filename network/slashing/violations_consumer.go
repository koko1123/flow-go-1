package slashing

import (
	"github.com/koko1123/flow-go-1/model/flow"
	network "github.com/koko1123/flow-go-1/network/channels"
	"github.com/koko1123/flow-go-1/network/message"
)

type ViolationsConsumer interface {
	// OnUnAuthorizedSenderError logs an error for unauthorized sender error
	OnUnAuthorizedSenderError(violation *Violation)

	// OnUnknownMsgTypeError logs an error for unknown message type error
	OnUnknownMsgTypeError(violation *Violation)

	// OnInvalidMsgError logs an error for messages that contained payloads that could not
	// be unmarshalled into the message type denoted by message code byte.
	OnInvalidMsgError(violation *Violation)

	// OnSenderEjectedError logs an error for sender ejected error
	OnSenderEjectedError(violation *Violation)

	// OnUnauthorizedUnicastOnChannel logs an error for messages unauthorized to be sent via unicast
	OnUnauthorizedUnicastOnChannel(violation *Violation)

	// OnUnexpectedError logs an error for unknown errors
	OnUnexpectedError(violation *Violation)
}

type Violation struct {
	Identity *flow.Identity
	PeerID   string
	MsgType  string
	Channel  network.Channel
	Protocol message.Protocol
	Err      error
}
