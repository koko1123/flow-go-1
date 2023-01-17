// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	context "context"

	flow "github.com/koko1123/flow-go-1/model/flow"
	entities "github.com/onflow/flow/protobuf/go/flow/entities"

	mock "github.com/stretchr/testify/mock"
)

// API is an autogenerated mock type for the API type
type API struct {
	mock.Mock
}

// GetExecutionDataByBlockID provides a mock function with given fields: ctx, blockID
func (_m *API) GetExecutionDataByBlockID(ctx context.Context, blockID flow.Identifier) (*entities.BlockExecutionData, error) {
	ret := _m.Called(ctx, blockID)

	var r0 *entities.BlockExecutionData
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier) *entities.BlockExecutionData); ok {
		r0 = rf(ctx, blockID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.BlockExecutionData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier) error); ok {
		r1 = rf(ctx, blockID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewAPI interface {
	mock.TestingT
	Cleanup(func())
}

// NewAPI creates a new instance of API. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAPI(t mockConstructorTestingTNewAPI) *API {
	mock := &API{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
