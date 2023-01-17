package benchmark

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/utils/unittest"
)

// TestNewStatsPusher is a simple test that starts and stops a pusher.
func TestNewStatsPusher(t *testing.T) {
	t.Parallel()

	p := NewStatsPusher(context.Background(), zerolog.Logger{}, "", "test", prometheus.DefaultGatherer)
	require.NotNil(t, p)

	unittest.RequireReturnsBefore(t, p.Stop, 1*time.Second, "pusher did not close in time")
}
