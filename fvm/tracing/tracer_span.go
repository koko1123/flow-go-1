package tracing

import (
	otelTrace "go.opentelemetry.io/otel/trace"

	"github.com/koko1123/flow-go-1/module"
	"github.com/koko1123/flow-go-1/module/trace"
)

type TracerSpan struct {
	module.Tracer

	otelTrace.Span

	ExtensiveTracing bool
}

func NewTracerSpan() TracerSpan {
	return TracerSpan{}
}

func NewMockTracerSpan() TracerSpan {
	return TracerSpan{
		Span: trace.NoopSpan,
	}
}

func (tracer TracerSpan) isTraceable() bool {
	return tracer.Tracer != nil && tracer.Span != nil
}

func (tracer TracerSpan) StartChildSpan(
	name trace.SpanName,
	options ...otelTrace.SpanStartOption,
) TracerSpan {
	child := trace.NoopSpan
	if tracer.isTraceable() {
		child = tracer.Tracer.StartSpanFromParent(tracer.Span, name, options...)
	}

	return TracerSpan{
		Tracer:           tracer.Tracer,
		Span:             child,
		ExtensiveTracing: tracer.ExtensiveTracing,
	}
}

func (tracer TracerSpan) StartExtensiveTracingChildSpan(
	name trace.SpanName,
	options ...otelTrace.SpanStartOption,
) TracerSpan {
	child := trace.NoopSpan
	if tracer.isTraceable() && tracer.ExtensiveTracing {
		child = tracer.Tracer.StartSpanFromParent(tracer.Span, name, options...)
	}

	return TracerSpan{
		Tracer:           tracer.Tracer,
		Span:             child,
		ExtensiveTracing: tracer.ExtensiveTracing,
	}
}
