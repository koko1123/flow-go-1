// Code generated by mockery v2.13.1. DO NOT EDIT.

package mockp2p

import (
	p2p "github.com/koko1123/flow-go-1/network/p2p"
	mock "github.com/stretchr/testify/mock"

	peer "github.com/libp2p/go-libp2p/core/peer"
)

// PubSubAdapter is an autogenerated mock type for the PubSubAdapter type
type PubSubAdapter struct {
	mock.Mock
}

// GetTopics provides a mock function with given fields:
func (_m *PubSubAdapter) GetTopics() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// Join provides a mock function with given fields: topic
func (_m *PubSubAdapter) Join(topic string) (p2p.Topic, error) {
	ret := _m.Called(topic)

	var r0 p2p.Topic
	if rf, ok := ret.Get(0).(func(string) p2p.Topic); ok {
		r0 = rf(topic)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(p2p.Topic)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(topic)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPeers provides a mock function with given fields: topic
func (_m *PubSubAdapter) ListPeers(topic string) []peer.ID {
	ret := _m.Called(topic)

	var r0 []peer.ID
	if rf, ok := ret.Get(0).(func(string) []peer.ID); ok {
		r0 = rf(topic)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]peer.ID)
		}
	}

	return r0
}

// RegisterTopicValidator provides a mock function with given fields: topic, topicValidator
func (_m *PubSubAdapter) RegisterTopicValidator(topic string, topicValidator p2p.TopicValidatorFunc) error {
	ret := _m.Called(topic, topicValidator)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, p2p.TopicValidatorFunc) error); ok {
		r0 = rf(topic, topicValidator)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UnregisterTopicValidator provides a mock function with given fields: topic
func (_m *PubSubAdapter) UnregisterTopicValidator(topic string) error {
	ret := _m.Called(topic)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewPubSubAdapter interface {
	mock.TestingT
	Cleanup(func())
}

// NewPubSubAdapter creates a new instance of PubSubAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPubSubAdapter(t mockConstructorTestingTNewPubSubAdapter) *PubSubAdapter {
	mock := &PubSubAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
