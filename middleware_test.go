package logz_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glassonion1/logz"
	"github.com/glassonion1/logz/internal/tracer"
)

func TestHTTPMiddleware(t *testing.T) {

	if err := logz.InitStdoutTracer(); err != nil {
		t.Fatalf("failed to init tracer: %v", err)
	}

	t.Run("Tests the middleware", func(t *testing.T) {
		defer func() {

		}()

		mux := http.NewServeMux()
		mux.Handle("/test1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logz.Infof(ctx, "write %s log", "info")

			traceID, spanID := tracer.TraceIDAndSpanID(ctx)

			if traceID == "00000000000000000000000000000000" {
				t.Error("trace id is zero value")
			}
			if spanID == "0000000000000000" {
				t.Error("span id is zero value")
			}
		}))

		mux.Handle("/test2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		req1 := httptest.NewRequest(http.MethodGet, "/test1", nil)
		rec1 := httptest.NewRecorder()
		mid.ServeHTTP(rec1, req1)

		req2 := httptest.NewRequest(http.MethodGet, "/test2", nil)
		rec2 := httptest.NewRecorder()
		mid.ServeHTTP(rec2, req2)
	})

}
