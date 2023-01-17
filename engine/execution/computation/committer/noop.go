package committer

import (
	"github.com/koko1123/flow-go-1/fvm/state"
	"github.com/koko1123/flow-go-1/ledger"
	"github.com/koko1123/flow-go-1/model/flow"
)

type NoopViewCommitter struct {
}

func NewNoopViewCommitter() *NoopViewCommitter {
	return &NoopViewCommitter{}
}

func (n NoopViewCommitter) CommitView(_ state.View, s flow.StateCommitment) (flow.StateCommitment, []byte, *ledger.TrieUpdate, error) {
	return s, nil, nil, nil
}
