package internal

import (
	"github.com/koko1123/flow-go-1/engine"
	"github.com/koko1123/flow-go-1/model/flow"
)

// MessageEntity is an internal data structure for storing messages in HeroQueue.
type MessageEntity struct {
	Msg engine.Message
	id  flow.Identifier
}

var _ flow.Entity = (*MessageEntity)(nil)

func NewMessageEntity(msg *engine.Message) MessageEntity {
	return MessageEntity{
		Msg: *msg,
		id:  identifierOfMessage(msg),
	}
}

func (m MessageEntity) ID() flow.Identifier {
	return m.id
}

func (m MessageEntity) Checksum() flow.Identifier {
	return m.id
}

func identifierOfMessage(msg *engine.Message) flow.Identifier {
	return flow.MakeID(msg)
}
