package tracer

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func TraceIDAndSpanID(ctx context.Context) (string, string) {
	spanCtx := trace.SpanContextFromContext(ctx)
	return spanCtx.TraceID.String(),
		spanCtx.SpanID.String()
}
