package mapfunc

import (
	"github.com/koko1123/flow-go-1/model/flow"
)

func WithWeight(weight uint64) flow.IdentityMapFunc {
	return func(identity flow.Identity) flow.Identity {
		identity.Weight = weight
		return identity
	}
}
