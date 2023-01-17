package topology

import (
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/network"
)

// EmptyTopology always returns an empty fanout list.
type EmptyTopology struct{}

var _ network.Topology = &EmptyTopology{}

func NewEmptyTopology() network.Topology {
	return &EmptyTopology{}
}

func (e EmptyTopology) Fanout(_ flow.IdentityList) flow.IdentityList {
	return flow.IdentityList{}
}
