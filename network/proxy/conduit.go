package proxy

import (
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/network"
)

// ProxyConduit is a special conduit which wraps the given conduit and replaces the target
// of every network send with the given target node.
type ProxyConduit struct {
	network.Conduit
	targetNodeID flow.Identifier
}

func (c *ProxyConduit) Publish(event interface{}, targetIDs ...flow.Identifier) error {
	return c.Conduit.Publish(event, c.targetNodeID)
}

func (c *ProxyConduit) Unicast(event interface{}, targetID flow.Identifier) error {
	return c.Conduit.Unicast(event, c.targetNodeID)
}

func (c *ProxyConduit) Multicast(event interface{}, num uint, targetIDs ...flow.Identifier) error {
	return c.Conduit.Multicast(event, 1, c.targetNodeID)
}
