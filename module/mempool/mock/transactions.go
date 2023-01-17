// Code generated by mockery v2.13.1. DO NOT EDIT.

package mempool

import (
	flow "github.com/koko1123/flow-go-1/model/flow"

	mock "github.com/stretchr/testify/mock"
)

// Transactions is an autogenerated mock type for the Transactions type
type Transactions struct {
	mock.Mock
}

// Add provides a mock function with given fields: tx
func (_m *Transactions) Add(tx *flow.TransactionBody) bool {
	ret := _m.Called(tx)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*flow.TransactionBody) bool); ok {
		r0 = rf(tx)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// All provides a mock function with given fields:
func (_m *Transactions) All() []*flow.TransactionBody {
	ret := _m.Called()

	var r0 []*flow.TransactionBody
	if rf, ok := ret.Get(0).(func() []*flow.TransactionBody); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*flow.TransactionBody)
		}
	}

	return r0
}

// ByID provides a mock function with given fields: txID
func (_m *Transactions) ByID(txID flow.Identifier) (*flow.TransactionBody, bool) {
	ret := _m.Called(txID)

	var r0 *flow.TransactionBody
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.TransactionBody); ok {
		r0 = rf(txID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.TransactionBody)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(flow.Identifier) bool); ok {
		r1 = rf(txID)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Clear provides a mock function with given fields:
func (_m *Transactions) Clear() {
	_m.Called()
}

// Has provides a mock function with given fields: txID
func (_m *Transactions) Has(txID flow.Identifier) bool {
	ret := _m.Called(txID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(flow.Identifier) bool); ok {
		r0 = rf(txID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Remove provides a mock function with given fields: txID
func (_m *Transactions) Remove(txID flow.Identifier) bool {
	ret := _m.Called(txID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(flow.Identifier) bool); ok {
		r0 = rf(txID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Size provides a mock function with given fields:
func (_m *Transactions) Size() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

type mockConstructorTestingTNewTransactions interface {
	mock.TestingT
	Cleanup(func())
}

// NewTransactions creates a new instance of Transactions. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTransactions(t mockConstructorTestingTNewTransactions) *Transactions {
	mock := &Transactions{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
