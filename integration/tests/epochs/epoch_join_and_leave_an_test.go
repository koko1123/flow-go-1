package epochs

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

func TestEpochJoinAndLeaveAN(t *testing.T) {
	unittest.SkipUnless(t, unittest.TEST_FLAKY, "epochs join/leave tests should be run on an machine with adequate resources")
	suite.Run(t, new(EpochJoinAndLeaveANSuite))
}

type EpochJoinAndLeaveANSuite struct {
	DynamicEpochTransitionSuite
}

// TestEpochJoinAndLeaveAN should update access nodes and assert healthy network conditions
// after the epoch transition completes. See health check function for details.
func (s *EpochJoinAndLeaveANSuite) TestEpochJoinAndLeaveAN() {
	s.runTestEpochJoinAndLeave(flow.RoleAccess, s.assertNetworkHealthyAfterANChange)
}
