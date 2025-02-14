package main

import (
	"math/rand"
	"time"

	"github.com/rs/zerolog"

	"github.com/koko1123/flow-go-1/module"
	"github.com/koko1123/flow-go-1/module/metrics"
	"github.com/koko1123/flow-go-1/module/metrics/example"
	"github.com/koko1123/flow-go-1/module/trace"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

// main runs a local tracer server on the machine and starts monitoring some metrics for sake of execution, which
// increases result approvals counter and checked chunks counter 100 times each
func main() {
	example.WithMetricsServer(func(logger zerolog.Logger) {
		tracer, err := trace.NewTracer(logger, "collection", "test", trace.SensitivityCaptureAll)
		if err != nil {
			panic(err)
		}
		collector := struct {
			*metrics.HotstuffCollector
			*metrics.ExecutionCollector
			*metrics.NetworkCollector
		}{
			HotstuffCollector:  metrics.NewHotstuffCollector("some_chain_id"),
			ExecutionCollector: metrics.NewExecutionCollector(tracer),
			NetworkCollector:   metrics.NewNetworkCollector(unittest.Logger()),
		}
		diskTotal := rand.Int63n(1024 * 1024 * 1024)
		for i := 0; i < 1000; i++ {
			blockID := unittest.BlockFixture().ID()
			collector.StartBlockReceivedToExecuted(blockID)

			duration := time.Duration(rand.Int31n(2000)) * time.Millisecond
			// adds a random delay for execution duration, between 0 and 2 seconds
			time.Sleep(duration)

			collector.ExecutionBlockExecuted(
				duration,
				module.ExecutionResultStats{
					ComputationUsed:      uint64(rand.Int63n(1e6)),
					MemoryUsed:           uint64(rand.Int63n(1e6)),
					EventCounts:          2,
					EventSize:            100,
					NumberOfCollections:  1,
					NumberOfTransactions: 1,
				})

			diskIncrease := rand.Int63n(1024 * 1024)
			diskTotal += diskIncrease
			collector.ExecutionStateStorageDiskTotal(diskTotal)
			collector.ExecutionStorageStateCommitment(diskIncrease)

			collector.FinishBlockReceivedToExecuted(blockID)
		}
	})
}
