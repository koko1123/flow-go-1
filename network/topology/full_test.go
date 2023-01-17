package topology_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/network/topology"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

// TestFullTopology checks that FullyConnectedTopology always returns the input list as the fanout.
func TestFullTopology(t *testing.T) {
	ids := unittest.IdentityListFixture(10)
	top := topology.NewFullyConnectedTopology()
	require.Equal(t, ids, top.Fanout(ids))
}
