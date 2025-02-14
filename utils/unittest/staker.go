package unittest

import (
	"github.com/koko1123/flow-go-1/model/flow"
)

type FixedStaker struct {
	Staked bool
}

func NewFixedStaker(initial bool) *FixedStaker {
	return &FixedStaker{
		Staked: initial,
	}
}

func (f *FixedStaker) AmIStakedAt(_ flow.Identifier) bool {
	return f.Staked
}
