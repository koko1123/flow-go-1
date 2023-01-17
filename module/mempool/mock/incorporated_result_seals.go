// Code generated by mockery v2.13.1. DO NOT EDIT.

package mempool

import (
	flow "github.com/koko1123/flow-go-1/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// IncorporatedResultSeals is an autogenerated mock type for the IncorporatedResultSeals type
type IncorporatedResultSeals struct {
	mock.Mock
}

// Add provides a mock function with given fields: irSeal
func (_m *IncorporatedResultSeals) Add(irSeal *flow.IncorporatedResultSeal) (bool, error) {
	ret := _m.Called(irSeal)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*flow.IncorporatedResultSeal) bool); ok {
		r0 = rf(irSeal)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*flow.IncorporatedResultSeal) error); ok {
		r1 = rf(irSeal)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// All provides a mock function with given fields:
func (_m *IncorporatedResultSeals) All() []*flow.IncorporatedResultSeal {
	ret := _m.Called()

	var r0 []*flow.IncorporatedResultSeal
	if rf, ok := ret.Get(0).(func() []*flow.IncorporatedResultSeal); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*flow.IncorporatedResultSeal)
		}
	}

	return r0
}

// ByID provides a mock function with given fields: _a0
func (_m *IncorporatedResultSeals) ByID(_a0 flow.Identifier) (*flow.IncorporatedResultSeal, bool) {
	ret := _m.Called(_a0)

	var r0 *flow.IncorporatedResultSeal
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.IncorporatedResultSeal); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.IncorporatedResultSeal)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(flow.Identifier) bool); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Clear provides a mock function with given fields:
func (_m *IncorporatedResultSeals) Clear() {
	_m.Called()
}

// Limit provides a mock function with given fields:
func (_m *IncorporatedResultSeals) Limit() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

// PruneUpToHeight provides a mock function with given fields: height
func (_m *IncorporatedResultSeals) PruneUpToHeight(height uint64) error {
	ret := _m.Called(height)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64) error); ok {
		r0 = rf(height)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Remove provides a mock function with given fields: incorporatedResultID
func (_m *IncorporatedResultSeals) Remove(incorporatedResultID flow.Identifier) bool {
	ret := _m.Called(incorporatedResultID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(flow.Identifier) bool); ok {
		r0 = rf(incorporatedResultID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Size provides a mock function with given fields:
func (_m *IncorporatedResultSeals) Size() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

type mockConstructorTestingTNewIncorporatedResultSeals interface {
	mock.TestingT
	Cleanup(func())
}

// NewIncorporatedResultSeals creates a new instance of IncorporatedResultSeals. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIncorporatedResultSeals(t mockConstructorTestingTNewIncorporatedResultSeals) *IncorporatedResultSeals {
	mock := &IncorporatedResultSeals{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
