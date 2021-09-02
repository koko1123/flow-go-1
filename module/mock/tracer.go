// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	context "context"

	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"

	opentracing "github.com/opentracing/opentracing-go"

	time "time"

	trace "github.com/onflow/flow-go/module/trace"
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

// EntityRootSpan provides a mock function with given fields: entityID, entityType, opts
func (_m *Tracer) EntityRootSpan(entityID flow.Identifier, entityType string, opts ...opentracing.StartSpanOption) opentracing.Span {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, entityID, entityType)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 opentracing.Span
	if rf, ok := ret.Get(0).(func(flow.Identifier, string, ...opentracing.StartSpanOption) opentracing.Span); ok {
		r0 = rf(entityID, entityType, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(opentracing.Span)
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

// RecordSpanFromParent provides a mock function with given fields: span, operationName, duration, logs, opts
func (_m *Tracer) RecordSpanFromParent(span opentracing.Span, operationName trace.SpanName, duration time.Duration, logs []opentracing.LogRecord, opts ...opentracing.StartSpanOption) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, span, operationName, duration, logs)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// StartBlockSpan provides a mock function with given fields: ctx, blockID, spanName, opts
func (_m *Tracer) StartBlockSpan(ctx context.Context, blockID flow.Identifier, spanName trace.SpanName, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, blockID, spanName)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 opentracing.Span
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier, trace.SpanName, ...opentracing.StartSpanOption) opentracing.Span); ok {
		r0 = rf(ctx, blockID, spanName, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(opentracing.Span)
		}
	}

	var r1 context.Context
	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier, trace.SpanName, ...opentracing.StartSpanOption) context.Context); ok {
		r1 = rf(ctx, blockID, spanName, opts...)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(context.Context)
		}
	}

	return r0, r1
}

// StartCollectionSpan provides a mock function with given fields: ctx, collectionID, spanName, opts
func (_m *Tracer) StartCollectionSpan(ctx context.Context, collectionID flow.Identifier, spanName trace.SpanName, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, collectionID, spanName)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 opentracing.Span
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier, trace.SpanName, ...opentracing.StartSpanOption) opentracing.Span); ok {
		r0 = rf(ctx, collectionID, spanName, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(opentracing.Span)
		}
	}

	var r1 context.Context
	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier, trace.SpanName, ...opentracing.StartSpanOption) context.Context); ok {
		r1 = rf(ctx, collectionID, spanName, opts...)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(context.Context)
		}
	}

	return r0, r1
}

// StartSpanFromContext provides a mock function with given fields: ctx, operationName, opts
func (_m *Tracer) StartSpanFromContext(ctx context.Context, operationName trace.SpanName, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, operationName)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 opentracing.Span
	if rf, ok := ret.Get(0).(func(context.Context, trace.SpanName, ...opentracing.StartSpanOption) opentracing.Span); ok {
		r0 = rf(ctx, operationName, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(opentracing.Span)
		}
	}

	var r1 context.Context
	if rf, ok := ret.Get(1).(func(context.Context, trace.SpanName, ...opentracing.StartSpanOption) context.Context); ok {
		r1 = rf(ctx, operationName, opts...)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(context.Context)
		}
	}

	return r0, r1
}

// StartSpanFromParent provides a mock function with given fields: span, operationName, opts
func (_m *Tracer) StartSpanFromParent(span opentracing.Span, operationName trace.SpanName, opts ...opentracing.StartSpanOption) opentracing.Span {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, span, operationName)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 opentracing.Span
	if rf, ok := ret.Get(0).(func(opentracing.Span, trace.SpanName, ...opentracing.StartSpanOption) opentracing.Span); ok {
		r0 = rf(span, operationName, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(opentracing.Span)
		}
	}

	return r0
}

// StartTransactionSpan provides a mock function with given fields: ctx, transactionID, spanName, opts
func (_m *Tracer) StartTransactionSpan(ctx context.Context, transactionID flow.Identifier, spanName trace.SpanName, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, transactionID, spanName)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 opentracing.Span
	if rf, ok := ret.Get(0).(func(context.Context, flow.Identifier, trace.SpanName, ...opentracing.StartSpanOption) opentracing.Span); ok {
		r0 = rf(ctx, transactionID, spanName, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(opentracing.Span)
		}
	}

	var r1 context.Context
	if rf, ok := ret.Get(1).(func(context.Context, flow.Identifier, trace.SpanName, ...opentracing.StartSpanOption) context.Context); ok {
		r1 = rf(ctx, transactionID, spanName, opts...)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(context.Context)
		}
	}

	return r0, r1
}

// WithSpanFromContext provides a mock function with given fields: ctx, operationName, f, opts
func (_m *Tracer) WithSpanFromContext(ctx context.Context, operationName trace.SpanName, f func(), opts ...opentracing.StartSpanOption) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, operationName, f)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}
