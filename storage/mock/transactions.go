// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	flow "github.com/koko1123/flow-go-1/model/flow"
	mock "github.com/stretchr/testify/mock"
)

// Transactions is an autogenerated mock type for the Transactions type
type Transactions struct {
	mock.Mock
}

// ByID provides a mock function with given fields: txID
func (_m *Transactions) ByID(txID flow.Identifier) (*flow.TransactionBody, error) {
	ret := _m.Called(txID)

	var r0 *flow.TransactionBody
	if rf, ok := ret.Get(0).(func(flow.Identifier) *flow.TransactionBody); ok {
		r0 = rf(txID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.TransactionBody)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(flow.Identifier) error); ok {
		r1 = rf(txID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: tx
func (_m *Transactions) Store(tx *flow.TransactionBody) error {
	ret := _m.Called(tx)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.TransactionBody) error); ok {
		r0 = rf(tx)
	} else {
		r0 = ret.Error(0)
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
