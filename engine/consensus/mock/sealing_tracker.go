// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	consensus "github.com/koko1123/flow-go-1/engine/consensus"
	flow "github.com/koko1123/flow-go-1/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// SealingTracker is an autogenerated mock type for the SealingTracker type
type SealingTracker struct {
	mock.Mock
}

// NewSealingObservation provides a mock function with given fields: finalizedBlock, seal, sealedBlock
func (_m *SealingTracker) NewSealingObservation(finalizedBlock *flow.Header, seal *flow.Seal, sealedBlock *flow.Header) consensus.SealingObservation {
	ret := _m.Called(finalizedBlock, seal, sealedBlock)

	var r0 consensus.SealingObservation
	if rf, ok := ret.Get(0).(func(*flow.Header, *flow.Seal, *flow.Header) consensus.SealingObservation); ok {
		r0 = rf(finalizedBlock, seal, sealedBlock)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(consensus.SealingObservation)
		}
	}

	return r0
}

type mockConstructorTestingTNewSealingTracker interface {
	mock.TestingT
	Cleanup(func())
}

// NewSealingTracker creates a new instance of SealingTracker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSealingTracker(t mockConstructorTestingTNewSealingTracker) *SealingTracker {
	mock := &SealingTracker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
