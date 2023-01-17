package validator

import (
	"github.com/koko1123/flow-go-1/module"
	"github.com/koko1123/flow-go-1/network"
)

var _ network.MessageValidator = (*OriginValidator)(nil)

// OriginValidator returns true if the sender of the message is among the set of identifiers
// returned by the given IdentifierProvider
type OriginValidator struct {
	idProvider module.IdentifierProvider
}

func NewOriginValidator(provider module.IdentifierProvider) network.MessageValidator {
	return &OriginValidator{provider}
}

func (v OriginValidator) Validate(msg network.IncomingMessageScope) bool {
	return v.idProvider.Identifiers().Contains(msg.OriginId())
}
