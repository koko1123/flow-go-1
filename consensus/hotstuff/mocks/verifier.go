// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	flow "github.com/koko1123/flow-go-1/model/flow"

	mock "github.com/stretchr/testify/mock"

	model "github.com/koko1123/flow-go-1/consensus/hotstuff/model"
)

// Verifier is an autogenerated mock type for the Verifier type
type Verifier struct {
	mock.Mock
}

// VerifyQC provides a mock function with given fields: signers, sigData, block
func (_m *Verifier) VerifyQC(signers flow.IdentityList, sigData []byte, block *model.Block) error {
	ret := _m.Called(signers, sigData, block)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.IdentityList, []byte, *model.Block) error); ok {
		r0 = rf(signers, sigData, block)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VerifyVote provides a mock function with given fields: voter, sigData, block
func (_m *Verifier) VerifyVote(voter *flow.Identity, sigData []byte, block *model.Block) error {
	ret := _m.Called(voter, sigData, block)

	var r0 error
	if rf, ok := ret.Get(0).(func(*flow.Identity, []byte, *model.Block) error); ok {
		r0 = rf(voter, sigData, block)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewVerifier interface {
	mock.TestingT
	Cleanup(func())
}

// NewVerifier creates a new instance of Verifier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewVerifier(t mockConstructorTestingTNewVerifier) *Verifier {
	mock := &Verifier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
