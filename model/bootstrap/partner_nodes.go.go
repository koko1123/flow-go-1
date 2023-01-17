package bootstrap

import (
	"github.com/koko1123/flow-go-1/model/encodable"
	"github.com/koko1123/flow-go-1/model/flow"
)

// PartnerNodeInfoPub represents public information about a partner/external
// node. It is identical to NodeInfoPub, but without weight information, as this
// is determined externally to the process that generates this information.
type PartnerNodeInfoPub struct {
	Role          flow.Role
	Address       string
	NodeID        flow.Identifier
	NetworkPubKey encodable.NetworkPubKey
	StakingPubKey encodable.StakingPubKey
}
