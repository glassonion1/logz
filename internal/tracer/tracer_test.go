package tracer_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/glassonion1/logz/internal/tracer"
)

func TestTracer(t *testing.T) {

	// Sets the trace exporter for stdout
	exporter, err := stdout.NewExporter([]stdout.Option{
		stdout.WithQuantiles([]float64{0.5, 0.9, 0.99}),
		stdout.WithPrettyPrint(),
	}...)
	if err != nil {
		t.Fatalf("failed to initialize stdout export pipeline: %v", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(bsp))
	defer func() { _ = tp.Shutdown(context.Background()) }()

	otel.SetTracerProvider(tp)

	t.Run("Tests the tracer", func(t *testing.T) {
		ctx := context.Background()
		tra := tracer.New("root tracer")
		childCtx, end := tra.Start(ctx, "test/trace")
		defer end()

		traceID, spanID := tracer.TraceIDAndSpanID(childCtx)

		if traceID == "00000000000000000000000000000000" {
			t.Error("trace id is zero value")
		}
		if spanID == "0000000000000000" {
			t.Error("span id is zero value")
		}
	})
}
