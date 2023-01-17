package hotstuff

import (
	"github.com/koko1123/flow-go-1/consensus/hotstuff/model"
)

// Signer is responsible for creating votes, proposals for a given block.
type Signer interface {
	// CreateProposal creates a proposal for the given block. No error returns
	// are expected during normal operations (incl. presence of byz. actors).
	CreateProposal(block *model.Block) (*model.Proposal, error)

	// CreateVote creates a vote for the given block. No error returns are
	// expected during normal operations (incl. presence of byz. actors).
	CreateVote(block *model.Block) (*model.Vote, error)
}
