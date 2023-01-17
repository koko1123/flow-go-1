// Code generated by mockery v2.13.1. DO NOT EDIT.

package component

import (
	component "github.com/koko1123/flow-go-1/module/component"
	mock "github.com/stretchr/testify/mock"
)

// ComponentManagerBuilder is an autogenerated mock type for the ComponentManagerBuilder type
type ComponentManagerBuilder struct {
	mock.Mock
}

// AddWorker provides a mock function with given fields: _a0
func (_m *ComponentManagerBuilder) AddWorker(_a0 component.ComponentWorker) component.ComponentManagerBuilder {
	ret := _m.Called(_a0)

	var r0 component.ComponentManagerBuilder
	if rf, ok := ret.Get(0).(func(component.ComponentWorker) component.ComponentManagerBuilder); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(component.ComponentManagerBuilder)
		}
	}

	return r0
}

// Build provides a mock function with given fields:
func (_m *ComponentManagerBuilder) Build() *component.ComponentManager {
	ret := _m.Called()

	var r0 *component.ComponentManager
	if rf, ok := ret.Get(0).(func() *component.ComponentManager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*component.ComponentManager)
		}
	}

	return r0
}

type mockConstructorTestingTNewComponentManagerBuilder interface {
	mock.TestingT
	Cleanup(func())
}

// NewComponentManagerBuilder creates a new instance of ComponentManagerBuilder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewComponentManagerBuilder(t mockConstructorTestingTNewComponentManagerBuilder) *ComponentManagerBuilder {
	mock := &ComponentManagerBuilder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
