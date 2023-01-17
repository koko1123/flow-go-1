package proxy

import (
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/network"
	"github.com/koko1123/flow-go-1/network/channels"
)

type ProxyNetwork struct {
	network.Network
	targetNodeID flow.Identifier
}

// NewProxyNetwork creates a new proxy network. All messages sent on this network are
// sent only to the node identified by the given target ID.
func NewProxyNetwork(net network.Network, targetNodeID flow.Identifier) *ProxyNetwork {
	return &ProxyNetwork{
		net,
		targetNodeID,
	}
}

// Register registers an engine with the proxy network.
func (n *ProxyNetwork) Register(channel channels.Channel, engine network.Engine) (network.Conduit, error) {
	con, err := n.Network.Register(channel, engine)

	if err != nil {
		return nil, err
	}

	proxyCon := ProxyConduit{
		con,
		n.targetNodeID,
	}

	return &proxyCon, nil
}
