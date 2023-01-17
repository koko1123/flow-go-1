package translator

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/model/flow/filter"
	"github.com/koko1123/flow-go-1/module"
	"github.com/koko1123/flow-go-1/network/p2p"
	"github.com/koko1123/flow-go-1/network/p2p/keyutils"
)

// IdentityProviderIDTranslator implements an `p2p.IDTranslator` which provides ID
// translation capabilities for an IdentityProvider.
type IdentityProviderIDTranslator struct {
	idProvider module.IdentityProvider
}

var _ p2p.IDTranslator = (*IdentityProviderIDTranslator)(nil)

func (t *IdentityProviderIDTranslator) GetFlowID(p peer.ID) (flow.Identifier, error) {
	key, err := p.ExtractPublicKey()
	if err != nil {
		return flow.ZeroID, err
	}
	flowKey, err := keyutils.FlowPublicKeyFromLibP2P(key)
	if err != nil {
		return flow.ZeroID, err
	}
	ids := t.idProvider.Identities(filter.HasNetworkingKey(flowKey))
	if len(ids) == 0 {
		return flow.ZeroID, fmt.Errorf("could not find identity corresponding to peer id %v", p.String())
	}
	return ids[0].NodeID, nil
}

func (t *IdentityProviderIDTranslator) GetPeerID(n flow.Identifier) (peer.ID, error) {
	ids := t.idProvider.Identities(filter.HasNodeID(n))
	if len(ids) == 0 {
		return "", fmt.Errorf("could not find identity with id %v", n.String())
	}
	key, err := keyutils.LibP2PPublicKeyFromFlow(ids[0].NetworkPubKey)
	if err != nil {
		return "", err
	}
	pid, err := peer.IDFromPublicKey(key)
	if err != nil {
		return "", err
	}
	return pid, nil
}

func NewIdentityProviderIDTranslator(provider module.IdentityProvider) *IdentityProviderIDTranslator {
	return &IdentityProviderIDTranslator{provider}
}
