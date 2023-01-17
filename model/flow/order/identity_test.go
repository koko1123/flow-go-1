// (c) 2019 Dapper Labs - ALL RIGHTS RESERVED

package order_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/model/flow/order"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

// Test the canonical ordering of identity and identifier match
func TestCanonicalOrderingMatch(t *testing.T) {
	identities := unittest.IdentityListFixture(100)
	require.Equal(t, identities.Sort(order.Canonical).NodeIDs(), identities.NodeIDs().Sort(order.IdentifierCanonical))
}
