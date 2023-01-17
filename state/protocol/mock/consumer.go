// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	flow "github.com/koko1123/flow-go-1/model/flow"
	mock "github.com/stretchr/testify/mock"
)

// Consumer is an autogenerated mock type for the Consumer type
type Consumer struct {
	mock.Mock
}

// BlockFinalized provides a mock function with given fields: block
func (_m *Consumer) BlockFinalized(block *flow.Header) {
	_m.Called(block)
}

// BlockProcessable provides a mock function with given fields: block
func (_m *Consumer) BlockProcessable(block *flow.Header) {
	_m.Called(block)
}

// EpochCommittedPhaseStarted provides a mock function with given fields: currentEpochCounter, first
func (_m *Consumer) EpochCommittedPhaseStarted(currentEpochCounter uint64, first *flow.Header) {
	_m.Called(currentEpochCounter, first)
}

// EpochSetupPhaseStarted provides a mock function with given fields: currentEpochCounter, first
func (_m *Consumer) EpochSetupPhaseStarted(currentEpochCounter uint64, first *flow.Header) {
	_m.Called(currentEpochCounter, first)
}

// EpochTransition provides a mock function with given fields: newEpochCounter, first
func (_m *Consumer) EpochTransition(newEpochCounter uint64, first *flow.Header) {
	_m.Called(newEpochCounter, first)
}

type mockConstructorTestingTNewConsumer interface {
	mock.TestingT
	Cleanup(func())
}

// NewConsumer creates a new instance of Consumer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConsumer(t mockConstructorTestingTNewConsumer) *Consumer {
	mock := &Consumer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
