package integration_test

import (
	"github.com/koko1123/flow-go-1/consensus/hotstuff/model"
	"github.com/koko1123/flow-go-1/consensus/hotstuff/notifications"
)

type CounterConsumer struct {
	notifications.NoopConsumer
	total     uint
	finalized func(uint)
}

func (c *CounterConsumer) OnFinalizedBlock(block *model.Block) {
	c.total++

	// notify stopper of total finalized
	c.finalized(c.total)
}
