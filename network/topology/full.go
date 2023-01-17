package topology

import (
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/network"
)

// FullyConnectedTopology returns all nodes as the fanout.
type FullyConnectedTopology struct{}

var _ network.Topology = &FullyConnectedTopology{}

func NewFullyConnectedTopology() network.Topology {
	return &FullyConnectedTopology{}
}

func (f FullyConnectedTopology) Fanout(ids flow.IdentityList) flow.IdentityList {
	return ids
}
