// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	flow "github.com/koko1123/flow-go-1/model/flow"
	mock "github.com/stretchr/testify/mock"

	transaction "github.com/koko1123/flow-go-1/storage/badger/transaction"
)

// EpochCommits is an autogenerated mock type for the EpochCommits type
type EpochCommits struct {
	mock.Mock
}

// ByID provides a mock function with given fields: _a0
func (_m *EpochCommits) ByID(_a0 flow.Identifier) (*flow.EpochCommit, error) {
	ret := _m.Called(_a0)

	var r0 *flow.EpochCommit
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.EpochCommit); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.EpochCommit)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StoreTx provides a mock function with given fields: commit
func (_m *EpochCommits) StoreTx(commit *flow.EpochCommit) func(*transaction.Tx) error {
	ret := _m.Called(commit)

	var r0 func(*transaction.Tx) error
	if rf, ok := ret.Get(0).(func(*flow.EpochCommit) func(*transaction.Tx) error); ok {
		r0 = rf(commit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(func(*transaction.Tx) error)
		}
	}

	return r0
}

type mockConstructorTestingTNewEpochCommits interface {
	mock.TestingT
	Cleanup(func())
}

// NewEpochCommits creates a new instance of EpochCommits. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEpochCommits(t mockConstructorTestingTNewEpochCommits) *EpochCommits {
	mock := &EpochCommits{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
