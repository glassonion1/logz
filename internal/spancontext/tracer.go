package spancontext

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// SpanContext is context value
type SpanContext struct {
	TraceID      string
	SpanID       string
	TraceSampled bool
}

// Extract extracts span context from context
func Extract(ctx context.Context) SpanContext {
	spanCtx := trace.SpanContextFromContext(ctx)
	return SpanContext{
		TraceID:      spanCtx.TraceID.String(),
		SpanID:       spanCtx.SpanID.String(),
		TraceSampled: spanCtx.IsSampled(),
	}
}
