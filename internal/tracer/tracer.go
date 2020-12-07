package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	tracer trace.Tracer
}

func New(label string) *Tracer {
	tracer := otel.Tracer(label)
	return &Tracer{
		tracer: tracer,
	}
}

func (t *Tracer) Start(ctx context.Context, label string) (context.Context, func()) {
	childCtx, span := t.tracer.Start(ctx, label)

	fn := func() {
		span.End()
	}

	return childCtx, fn
}

func TraceIDAndSpanID(ctx context.Context) (string, string) {
	spanCtx := trace.SpanContextFromContext(ctx)
	return spanCtx.TraceID.String(),
		spanCtx.SpanID.String()
}
