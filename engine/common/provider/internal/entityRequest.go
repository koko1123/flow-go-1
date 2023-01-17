package internal

import "github.com/koko1123/flow-go-1/model/flow"

type EntityRequest struct {
	OriginId  flow.Identifier
	EntityIds []flow.Identifier
	Nonce     uint64
}
