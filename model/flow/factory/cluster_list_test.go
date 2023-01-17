package factory_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/model/flow/factory"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

// NewClusterList assumes the input assignments are sorted, and fail if not.
// This tests verifies that NewClusterList has implemented the check on the assumption.
func TestNewClusterListFail(t *testing.T) {
	identities := unittest.IdentityListFixture(100, unittest.WithRole(flow.RoleCollection))
	assignments := unittest.ClusterAssignment(10, identities)

	tmp := assignments[1][0]
	assignments[1][0] = assignments[1][1]
	assignments[1][1] = tmp

	_, err := factory.NewClusterList(assignments, identities)
	require.Error(t, err)
}
