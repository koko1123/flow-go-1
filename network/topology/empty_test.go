package topology_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/network/topology"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

// TestEmptyTopology checks that EmptyTopology always creates an empty list of fanout.
func TestEmptyTopology(t *testing.T) {
	ids := unittest.IdentityListFixture(10)
	top := topology.NewEmptyTopology()
	require.Empty(t, top.Fanout(ids))
}
