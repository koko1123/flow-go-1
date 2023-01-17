package epochmgr

import (
	"github.com/koko1123/flow-go-1/consensus/hotstuff"
	"github.com/koko1123/flow-go-1/module"
	"github.com/koko1123/flow-go-1/network"
	"github.com/koko1123/flow-go-1/state/cluster"
	"github.com/koko1123/flow-go-1/state/protocol"
)

// EpochComponentsFactory is responsible for creating epoch-scoped components
// managed by the epoch manager engine for the given epoch.
type EpochComponentsFactory interface {

	// Create sets up and instantiates all dependencies for the epoch. It may
	// be used either for an ongoing epoch (for example, after a restart) or
	// for an epoch that will start soon. It is safe to call multiple times for
	// a given epoch counter.
	//
	// Must return ErrNotAuthorizedForEpoch if this node is not authorized in the epoch.
	Create(epoch protocol.Epoch) (
		state cluster.State,
		proposal network.Engine,
		sync network.Engine,
		hotstuff module.HotStuff,
		voteAggregator hotstuff.VoteAggregator,
		err error,
	)
}
