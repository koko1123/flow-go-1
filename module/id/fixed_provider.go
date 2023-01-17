package id

import (
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module"
	"github.com/koko1123/flow-go-1/network/p2p/keyutils"
)

// FixedIdentifierProvider implements an IdentifierProvider which provides a fixed list
// of identifiers.
type FixedIdentifierProvider struct {
	identifiers flow.IdentifierList
}

func NewFixedIdentifierProvider(identifiers flow.IdentifierList) *FixedIdentifierProvider {
	return &FixedIdentifierProvider{identifiers}
}

func (p *FixedIdentifierProvider) Identifiers() flow.IdentifierList {
	return p.identifiers
}

// FixedIdentityProvider implements an IdentityProvider which provides a fixed list
// of identities.
type FixedIdentityProvider struct {
	identities flow.IdentityList
}

var _ module.IdentityProvider = (*FixedIdentityProvider)(nil)

func NewFixedIdentityProvider(identities flow.IdentityList) *FixedIdentityProvider {
	return &FixedIdentityProvider{identities}
}

func (p *FixedIdentityProvider) Identities(filter flow.IdentityFilter) flow.IdentityList {
	return p.identities.Filter(filter)
}

func (p *FixedIdentityProvider) ByNodeID(flowID flow.Identifier) (*flow.Identity, bool) {
	for _, v := range p.identities {
		if v.ID() == flowID {
			return v, true
		}
	}
	return nil, false
}

func (p *FixedIdentityProvider) ByPeerID(peerID peer.ID) (*flow.Identity, bool) {
	for _, v := range p.identities {
		if id, err := keyutils.PeerIDFromFlowPublicKey(v.NetworkPubKey); err == nil {
			if id == peerID {
				return v, true
			}
		}

	}
	return nil, false

}
