package spancontext

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type SpanContext struct {
	TraceID string
	SpanID  string
}

func Extract(ctx context.Context) SpanContext {
	spanCtx := trace.SpanContextFromContext(ctx)
	return SpanContext{
		TraceID: spanCtx.TraceID.String(),
		SpanID:  spanCtx.SpanID.String(),
	}
}
