package consensus

import (
	"github.com/koko1123/flow-go-1/consensus/hotstuff/model"
	"github.com/koko1123/flow-go-1/consensus/hotstuff/notifications"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module"
)

// MetricsConsumer is a consumer that subscribes to hotstuff events and
// collects metrics data when certain events trigger.
// It depends on Metrics module to report metrics data.
type MetricsConsumer struct {
	// inherit from noop consumer in order to satisfy the full interface
	notifications.NoopConsumer
	metrics module.HotstuffMetrics
}

func NewMetricsConsumer(metrics module.HotstuffMetrics) *MetricsConsumer {
	return &MetricsConsumer{
		metrics: metrics,
	}
}

func (c *MetricsConsumer) OnEnteringView(view uint64, leader flow.Identifier) {
	c.metrics.SetCurView(view)
}

func (c *MetricsConsumer) OnQcIncorporated(qc *flow.QuorumCertificate) {
	c.metrics.SetQCView(qc.View)
}

func (c *MetricsConsumer) OnQcTriggeredViewChange(qc *flow.QuorumCertificate, newView uint64) {
	c.metrics.CountSkipped()
}

func (c *MetricsConsumer) OnReachedTimeout(info *model.TimerInfo) {
	c.metrics.CountTimeout()
}

func (c *MetricsConsumer) OnStartingTimeout(info *model.TimerInfo) {
	c.metrics.SetTimeout(info.Duration)
}
