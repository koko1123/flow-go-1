// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	flow "github.com/koko1123/flow-go-1/model/flow"

	mock "github.com/stretchr/testify/mock"

	model "github.com/koko1123/flow-go-1/consensus/hotstuff/model"

	time "time"
)

// EventHandlerV2 is an autogenerated mock type for the EventHandlerV2 type
type EventHandlerV2 struct {
	mock.Mock
}

// OnLocalTimeout provides a mock function with given fields:
func (_m *EventHandlerV2) OnLocalTimeout() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OnQCConstructed provides a mock function with given fields: qc
func (_m *EventHandlerV2) OnQCConstructed(qc *flow.QuorumCertificate) error {
	ret := _m.Called(qc)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.QuorumCertificate) error); ok {
		r0 = rf(qc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OnReceiveProposal provides a mock function with given fields: proposal
func (_m *EventHandlerV2) OnReceiveProposal(proposal *model.Proposal) error {
	ret := _m.Called(proposal)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Proposal) error); ok {
		r0 = rf(proposal)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields:
func (_m *EventHandlerV2) Start() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TimeoutChannel provides a mock function with given fields:
func (_m *EventHandlerV2) TimeoutChannel() <-chan time.Time {
	ret := _m.Called()

	var r0 <-chan time.Time
	if rf, ok := ret.Get(0).(func() <-chan time.Time); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan time.Time)
		}
	}

	return r0
}
