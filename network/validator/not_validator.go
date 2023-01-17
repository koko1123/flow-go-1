package validator

import (
	"github.com/koko1123/flow-go-1/network"
)

var _ network.MessageValidator = (*NotValidator)(nil)

// NotValidator returns the opposite result of the given validator for the Validate call
type NotValidator struct {
	validator network.MessageValidator
}

func NewNotValidator(validator network.MessageValidator) network.MessageValidator {
	return &NotValidator{
		validator: validator,
	}
}

func (n NotValidator) Validate(msg network.IncomingMessageScope) bool {
	return !n.validator.Validate(msg)
}
