// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	context "context"

	flow "github.com/koko1123/flow-go-1/model/flow"
	execution_data "github.com/koko1123/flow-go-1/module/executiondatasync/execution_data"

	mock "github.com/stretchr/testify/mock"
)

// ExecutionDataStore is an autogenerated mock type for the ExecutionDataStore type
type ExecutionDataStore struct {
	mock.Mock
}

// AddExecutionData provides a mock function with given fields: ctx, executionData
func (_m *ExecutionDataStore) AddExecutionData(ctx context.Context, executionData *execution_data.BlockExecutionData) (flow.Identifier, error) {
	ret := _m.Called(ctx, executionData)

	var r0 flow.Identifier
	if rf, ok := ret.Get(0).(func(context.Context, *execution_data.BlockExecutionData) flow.Identifier); ok {
		r0 = rf(ctx, executionData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flow.Identifier)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *execution_data.BlockExecutionData) error); ok {
		r1 = rf(ctx, executionData)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetExecutionData provides a mock function with given fields: ctx, rootID
func (_m *ExecutionDataStore) GetExecutionData(ctx context.Context, rootID flow.Identifier) (*execution_data.BlockExecutionData, error) {
	ret := _m.Called(ctx, rootID)

	var r0 *execution_data.BlockExecutionData
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) *execution_data.BlockExecutionData); ok {
		r0 = rf(ctx, rootID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*execution_data.BlockExecutionData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(ctx, rootID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewExecutionDataStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewExecutionDataStore creates a new instance of ExecutionDataStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewExecutionDataStore(t mockConstructorTestingTNewExecutionDataStore) *ExecutionDataStore {
	mock := &ExecutionDataStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
