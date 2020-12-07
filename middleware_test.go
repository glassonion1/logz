package logz_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/glassonion1/logz"
	"github.com/glassonion1/logz/internal/tracer"
)

func TestHTTPMiddleware(t *testing.T) {

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

	t.Run("Tests the middleware", func(t *testing.T) {
		defer func() {

		}()

		mux := http.NewServeMux()
		mux.Handle("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logz.Infof(ctx, "write %s log", "info")

			traceID, spanID := tracer.TraceIDAndSpanID(ctx)

			if traceID == "00000000000000000000000000000000" {
				t.Error("trace id is zero value")
			}
			if spanID == "0000000000000000" {
				t.Error("span id is zero value")
			}

			// nested function
			func(ctx context.Context) {
				logz.Infof(ctx, "write %s nested log", "info")

				childTraceID, childSpanID := tracer.TraceIDAndSpanID(ctx)
				if childTraceID != traceID {
					t.Error("trace and child trace id are not equal")
				}
				if childSpanID != spanID {
					t.Error("span and child span id are not equal")
				}
			}(ctx)
		}))

		mid := logz.HTTPMiddleware("test/component")(mux)
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		mid.ServeHTTP(rec, req)
	})

}
