// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	flow "github.com/koko1123/flow-go-1/model/flow"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// ComplianceMetrics is an autogenerated mock type for the ComplianceMetrics type
type ComplianceMetrics struct {
	mock.Mock
}

// BlockFinalized provides a mock function with given fields: _a0
func (_m *ComplianceMetrics) BlockFinalized(_a0 *flow.Block) {
	_m.Called(_a0)
}

// BlockProposalDuration provides a mock function with given fields: duration
func (_m *ComplianceMetrics) BlockProposalDuration(duration time.Duration) {
	_m.Called(duration)
}

// BlockSealed provides a mock function with given fields: _a0
func (_m *ComplianceMetrics) BlockSealed(_a0 *flow.Block) {
	_m.Called(_a0)
}

// CommittedEpochFinalView provides a mock function with given fields: view
func (_m *ComplianceMetrics) CommittedEpochFinalView(view uint64) {
	_m.Called(view)
}

// CurrentDKGPhase1FinalView provides a mock function with given fields: view
func (_m *ComplianceMetrics) CurrentDKGPhase1FinalView(view uint64) {
	_m.Called(view)
}

// CurrentDKGPhase2FinalView provides a mock function with given fields: view
func (_m *ComplianceMetrics) CurrentDKGPhase2FinalView(view uint64) {
	_m.Called(view)
}

// CurrentDKGPhase3FinalView provides a mock function with given fields: view
func (_m *ComplianceMetrics) CurrentDKGPhase3FinalView(view uint64) {
	_m.Called(view)
}

// CurrentEpochCounter provides a mock function with given fields: counter
func (_m *ComplianceMetrics) CurrentEpochCounter(counter uint64) {
	_m.Called(counter)
}

// CurrentEpochFinalView provides a mock function with given fields: view
func (_m *ComplianceMetrics) CurrentEpochFinalView(view uint64) {
	_m.Called(view)
}

// CurrentEpochPhase provides a mock function with given fields: phase
func (_m *ComplianceMetrics) CurrentEpochPhase(phase flow.EpochPhase) {
	_m.Called(phase)
}

// EpochEmergencyFallbackTriggered provides a mock function with given fields:
func (_m *ComplianceMetrics) EpochEmergencyFallbackTriggered() {
	_m.Called()
}

// FinalizedHeight provides a mock function with given fields: height
func (_m *ComplianceMetrics) FinalizedHeight(height uint64) {
	_m.Called(height)
}

// SealedHeight provides a mock function with given fields: height
func (_m *ComplianceMetrics) SealedHeight(height uint64) {
	_m.Called(height)
}

type mockConstructorTestingTNewComplianceMetrics interface {
	mock.TestingT
	Cleanup(func())
}

// NewComplianceMetrics creates a new instance of ComplianceMetrics. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewComplianceMetrics(t mockConstructorTestingTNewComplianceMetrics) *ComplianceMetrics {
	mock := &ComplianceMetrics{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
