package main

import (
	"github.com/rs/zerolog"

	"github.com/koko1123/flow-go-1/consensus/hotstuff/notifications"
	"github.com/koko1123/flow-go-1/consensus/hotstuff/notifications/pubsub"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module"
	metricsconsumer "github.com/koko1123/flow-go-1/module/metrics/hotstuff"
)

func createNotifier(log zerolog.Logger, metrics module.HotstuffMetrics, tracer module.Tracer, chain flow.ChainID,
) *pubsub.Distributor {
	telemetryConsumer := notifications.NewTelemetryConsumer(log, chain)
	metricsConsumer := metricsconsumer.NewMetricsConsumer(metrics)
	dis := pubsub.NewDistributor()
	dis.AddConsumer(telemetryConsumer)
	dis.AddConsumer(metricsConsumer)
	return dis
}
