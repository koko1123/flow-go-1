// Code generated by mockery v2.13.1. DO NOT EDIT.

package mempool

import (
	flow "github.com/koko1123/flow-go-1/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// BackData is an autogenerated mock type for the BackData type
type BackData struct {
	mock.Mock
}

// Add provides a mock function with given fields: entityID, entity
func (_m *BackData) Add(entityID flow.Identifier, entity flow.Entity) bool {
	ret := _m.Called(entityID, entity)

	var r0 bool
	if rf, ok := ret.Get(0).(func(flow.Identifier, flow.Entity) bool); ok {
		r0 = rf(entityID, entity)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Adjust provides a mock function with given fields: entityID, f
func (_m *BackData) Adjust(entityID flow.Identifier, f func(flow.Entity) flow.Entity) (flow.Entity, bool) {
	ret := _m.Called(entityID, f)

	var r0 flow.Entity
	if rf, ok := ret.Get(0).(func(flow.Identifier, func(flow.Entity) flow.Entity) flow.Entity); ok {
		r0 = rf(entityID, f)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flow.Entity)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(flow.Identifier, func(flow.Entity) flow.Entity) bool); ok {
		r1 = rf(entityID, f)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// All provides a mock function with given fields:
func (_m *BackData) All() map[flow.Identifier]flow.Entity {
	ret := _m.Called()

	var r0 map[flow.Identifier]flow.Entity
	if rf, ok := ret.Get(0).(func() map[flow.Identifier]flow.Entity); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[flow.Identifier]flow.Entity)
		}
	}

	return r0
}

// ByID provides a mock function with given fields: entityID
func (_m *BackData) ByID(entityID flow.Identifier) (flow.Entity, bool) {
	ret := _m.Called(entityID)

	var r0 flow.Entity
	if rf, ok := ret.Get(0).(func(flow.Identifier) flow.Entity); ok {
		r0 = rf(entityID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flow.Entity)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(flow.Identifier) bool); ok {
		r1 = rf(entityID)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Clear provides a mock function with given fields:
func (_m *BackData) Clear() {
	_m.Called()
}

// Entities provides a mock function with given fields:
func (_m *BackData) Entities() []flow.Entity {
	ret := _m.Called()

	var r0 []flow.Entity
	if rf, ok := ret.Get(0).(func() []flow.Entity); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]flow.Entity)
		}
	}

	return r0
}

// Has provides a mock function with given fields: entityID
func (_m *BackData) Has(entityID flow.Identifier) bool {
	ret := _m.Called(entityID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(flow.Identifier) bool); ok {
		r0 = rf(entityID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Identifiers provides a mock function with given fields:
func (_m *BackData) Identifiers() flow.IdentifierList {
	ret := _m.Called()

	var r0 flow.IdentifierList
	if rf, ok := ret.Get(0).(func() flow.IdentifierList); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flow.IdentifierList)
		}
	}

	return r0
}

// Remove provides a mock function with given fields: entityID
func (_m *BackData) Remove(entityID flow.Identifier) (flow.Entity, bool) {
	ret := _m.Called(entityID)

	var r0 flow.Entity
	if rf, ok := ret.Get(0).(func(flow.Identifier) flow.Entity); ok {
		r0 = rf(entityID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flow.Entity)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(flow.Identifier) bool); ok {
		r1 = rf(entityID)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Size provides a mock function with given fields:
func (_m *BackData) Size() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

type mockConstructorTestingTNewBackData interface {
	mock.TestingT
	Cleanup(func())
}

// NewBackData creates a new instance of BackData. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBackData(t mockConstructorTestingTNewBackData) *BackData {
	mock := &BackData{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
