package scoring

import (
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module"
)

// HasValidFlowIdentity checks if the peer has a valid Flow identity.
func HasValidFlowIdentity(idProvider module.IdentityProvider, pid peer.ID) (*flow.Identity, error) {
	flowId, ok := idProvider.ByPeerID(pid)
	if !ok {
		return nil, NewInvalidPeerIDError(pid, PeerIdStatusUnknown)
	}

	if flowId.Ejected {
		return nil, NewInvalidPeerIDError(pid, PeerIdStatusEjected)
	}

	return flowId, nil
}
