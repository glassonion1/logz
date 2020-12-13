package spancontext

import (
	"context"

	"go.opencensus.io/trace"
)

type SpanContext struct {
	TraceID      string
	SpanID       string
	TraceSampled bool
}

func Extract(ctx context.Context) SpanContext {
	spanCtx := trace.FromContext(ctx).SpanContext()
	return SpanContext{
		TraceID:      spanCtx.TraceID.String(),
		SpanID:       spanCtx.SpanID.String(),
		TraceSampled: spanCtx.TraceOptions.IsSampled(),
	}
}
