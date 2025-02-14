// Code generated by mockery v2.13.1. DO NOT EDIT.

package mock

import (
	context "context"

	flow "github.com/koko1123/flow-go-1/model/flow"
	mock "github.com/stretchr/testify/mock"

	oteltrace "go.opentelemetry.io/otel/trace"

	trace "github.com/koko1123/flow-go-1/module/trace"
)

// Tracer is an autogenerated mock type for the Tracer type
type Tracer struct {
	mock.Mock
}

// Done provides a mock function with given fields:
func (_m *Tracer) Done() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// Ready provides a mock function with given fields:
func (_m *Tracer) Ready() <-chan struct{} {
	ret := _m.Called()

	var r0 <-chan struct{}
	if rf, ok := ret.Get(0).(func() <-chan struct{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan struct{})
		}
	}

	return r0
}

// StartBlockSpan provides a mock function with given fields: ctx, blockID, spanName, opts
func (_m *Tracer) StartBlockSpan(ctx context.Context, blockID flow.Identifier, spanName trace.SpanName, opts ...oteltrace.SpanStartOption) (oteltrace.Span, context.Context) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, blockID, spanName)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 oteltrace.Span
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier, trace.SpanName, ...oteltrace.SpanStartOption) oteltrace.Span); ok {
		r0 = rf(ctx, blockID, spanName, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(oteltrace.Span)
		}
	}

	var r1 context.Context
	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier, trace.SpanName, ...oteltrace.SpanStartOption) context.Context); ok {
		r1 = rf(ctx, blockID, spanName, opts...)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(context.Context)
		}
	}

	return r0, r1
}

// StartCollectionSpan provides a mock function with given fields: ctx, collectionID, spanName, opts
func (_m *Tracer) StartCollectionSpan(ctx context.Context, collectionID flow.Identifier, spanName trace.SpanName, opts ...oteltrace.SpanStartOption) (oteltrace.Span, context.Context) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, collectionID, spanName)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 oteltrace.Span
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier, trace.SpanName, ...oteltrace.SpanStartOption) oteltrace.Span); ok {
		r0 = rf(ctx, collectionID, spanName, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(oteltrace.Span)
		}
	}

	var r1 context.Context
	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier, trace.SpanName, ...oteltrace.SpanStartOption) context.Context); ok {
		r1 = rf(ctx, collectionID, spanName, opts...)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(context.Context)
		}
	}

	return r0, r1
}

// StartSpanFromContext provides a mock function with given fields: ctx, operationName, opts
func (_m *Tracer) StartSpanFromContext(ctx context.Context, operationName trace.SpanName, opts ...oteltrace.SpanStartOption) (oteltrace.Span, context.Context) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, operationName)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 oteltrace.Span
	if rf, ok := ret.Get(0).(func(context.Context, trace.SpanName, ...oteltrace.SpanStartOption) oteltrace.Span); ok {
		r0 = rf(ctx, operationName, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(oteltrace.Span)
		}
	}

	var r1 context.Context
	if rf, ok := ret.Get(1).(func(context.Context, trace.SpanName, ...oteltrace.SpanStartOption) context.Context); ok {
		r1 = rf(ctx, operationName, opts...)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(context.Context)
		}
	}

	return r0, r1
}

// StartSpanFromParent provides a mock function with given fields: parentSpan, operationName, opts
func (_m *Tracer) StartSpanFromParent(parentSpan oteltrace.Span, operationName trace.SpanName, opts ...oteltrace.SpanStartOption) oteltrace.Span {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, parentSpan, operationName)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 oteltrace.Span
	if rf, ok := ret.Get(0).(func(oteltrace.Span, trace.SpanName, ...oteltrace.SpanStartOption) oteltrace.Span); ok {
		r0 = rf(parentSpan, operationName, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(oteltrace.Span)
		}
	}

	return r0
}

// StartTransactionSpan provides a mock function with given fields: ctx, transactionID, spanName, opts
func (_m *Tracer) StartTransactionSpan(ctx context.Context, transactionID flow.Identifier, spanName trace.SpanName, opts ...oteltrace.SpanStartOption) (oteltrace.Span, context.Context) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, transactionID, spanName)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 oteltrace.Span
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier, trace.SpanName, ...oteltrace.SpanStartOption) oteltrace.Span); ok {
		r0 = rf(ctx, transactionID, spanName, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(oteltrace.Span)
		}
	}

	var r1 context.Context
	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier, trace.SpanName, ...oteltrace.SpanStartOption) context.Context); ok {
		r1 = rf(ctx, transactionID, spanName, opts...)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(context.Context)
		}
	}

	return r0, r1
}

// WithSpanFromContext provides a mock function with given fields: ctx, operationName, f, opts
func (_m *Tracer) WithSpanFromContext(ctx context.Context, operationName trace.SpanName, f func(), opts ...oteltrace.SpanStartOption) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, operationName, f)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

type mockConstructorTestingTNewTracer interface {
	mock.TestingT
	Cleanup(func())
}

// NewTracer creates a new instance of Tracer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTracer(t mockConstructorTestingTNewTracer) *Tracer {
	mock := &Tracer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
