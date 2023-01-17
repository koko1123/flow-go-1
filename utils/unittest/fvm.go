package unittest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koko1123/flow-go-1/fvm/systemcontracts"
	"github.com/koko1123/flow-go-1/model/flow"
)

func IsServiceEvent(event flow.Event, chainID flow.ChainID) bool {
	serviceEvents, _ := systemcontracts.ServiceEventsForChain(chainID)
	for _, serviceEvent := range serviceEvents.All() {
		if serviceEvent.EventType() == event.Type {
			return true
		}
	}
	return false
}

// EnsureEventsIndexSeq checks if values of given event index sequence are monotonically increasing.
func EnsureEventsIndexSeq(t *testing.T, events []flow.Event, chainID flow.ChainID) {
	expectedEventIndex := uint32(0)
	for _, event := range events {
		require.Equal(t, expectedEventIndex, event.EventIndex)
		if IsServiceEvent(event, chainID) {
			// TODO: we will need to address the double counting issue for service events.
			//		 https://github.com/koko1123/flow-go-1/issues/3393
			expectedEventIndex += 2
		} else {
			expectedEventIndex++
		}
	}
}
