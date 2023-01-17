package hotstuff

import (
	"github.com/koko1123/flow-go-1/consensus/hotstuff/model"
)

// FollowerLogic runs a state machine to process proposals
type FollowerLogic interface {
	// FinalizedBlock returns the latest finalized block
	FinalizedBlock() *model.Block

	// AddBlock processes a block proposal
	AddBlock(proposal *model.Proposal) error
}
