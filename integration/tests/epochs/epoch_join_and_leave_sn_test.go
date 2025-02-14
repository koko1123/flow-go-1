package epochs

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

func TestEpochJoinAndLeaveSN(t *testing.T) {
	unittest.SkipUnless(t, unittest.TEST_FLAKY, "epochs join/leave tests should be run on an machine with adequate resources")
	suite.Run(t, new(EpochJoinAndLeaveSNSuite))
}

type EpochJoinAndLeaveSNSuite struct {
	DynamicEpochTransitionSuite
}

// TestEpochJoinAndLeaveSN should update consensus nodes and assert healthy network conditions
// after the epoch transition completes. See health check function for details.
func (s *EpochJoinAndLeaveSNSuite) TestEpochJoinAndLeaveSN() {
	s.runTestEpochJoinAndLeave(flow.RoleConsensus, s.assertNetworkHealthyAfterSNChange)
}
