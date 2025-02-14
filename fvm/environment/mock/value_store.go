// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	atree "github.com/onflow/atree"

	mock "github.com/stretchr/testify/mock"
)

// ValueStore is an autogenerated mock type for the ValueStore type
type ValueStore struct {
	mock.Mock
}

// AllocateStorageIndex provides a mock function with given fields: owner
func (_m *ValueStore) AllocateStorageIndex(owner []byte) (atree.StorageIndex, error) {
	ret := _m.Called(owner)

	var r0 atree.StorageIndex
	if rf, ok := ret.Get(0).(func([]byte) atree.StorageIndex); ok {
		r0 = rf(owner)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(atree.StorageIndex)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(owner)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetValue provides a mock function with given fields: owner, key
func (_m *ValueStore) GetValue(owner []byte, key []byte) ([]byte, error) {
	ret := _m.Called(owner, key)

	var r0 []byte
	if rf, ok := ret.Get(0).(func([]byte, []byte) []byte); ok {
		r0 = rf(owner, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte, []byte) error); ok {
		r1 = rf(owner, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetValue provides a mock function with given fields: owner, key, value
func (_m *ValueStore) SetValue(owner []byte, key []byte, value []byte) error {
	ret := _m.Called(owner, key, value)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, []byte, []byte) error); ok {
		r0 = rf(owner, key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValueExists provides a mock function with given fields: owner, key
func (_m *ValueStore) ValueExists(owner []byte, key []byte) (bool, error) {
	ret := _m.Called(owner, key)

	var r0 bool
	if rf, ok := ret.Get(0).(func([]byte, []byte) bool); ok {
		r0 = rf(owner, key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte, []byte) error); ok {
		r1 = rf(owner, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewValueStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewValueStore creates a new instance of ValueStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewValueStore(t mockConstructorTestingTNewValueStore) *ValueStore {
	mock := &ValueStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
