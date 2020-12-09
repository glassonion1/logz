package logz_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glassonion1/logz"
	"github.com/glassonion1/logz/internal/tracer"
	"github.com/google/go-cmp/cmp"
	"github.com/googleinterns/cloud-operations-api-mock/cloudmock"
	"google.golang.org/api/option"
)

func TestHTTPMiddlewareWithStdoutTracer(t *testing.T) {

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

func TestHTTPMiddlewareWithCloudTracer(t *testing.T) {

	mock := cloudmock.NewCloudMock()
	defer mock.Shutdown()
	clientOpts := []option.ClientOption{
		option.WithGRPCConn(mock.ClientConn()),
	}

	if err := logz.InitCloudTracer(clientOpts...); err != nil {
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

func TestHTTPMiddlewareRemoteParent(t *testing.T) {

	mock := cloudmock.NewCloudMock()
	defer mock.Shutdown()
	clientOpts := []option.ClientOption{
		option.WithGRPCConn(mock.ClientConn()),
	}

	if err := logz.InitCloudTracer(clientOpts...); err != nil {
		t.Fatalf("failed to init tracer: %v", err)
	}

	t.Run("Tests the middleware with remote parent", func(t *testing.T) {
		defer func() {

		}()

		mux := http.NewServeMux()
		mux.Handle("/test1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logz.Infof(ctx, "write %s log", "info")

			traceID, spanID := tracer.TraceIDAndSpanID(ctx)

			if diff := cmp.Diff(traceID, "a0d3eee13de6a4bbcf291eb444b94f28"); diff != "" {
				t.Errorf("remote and current trace id are missmatch: %v", diff)
			}
			if spanID == "0000000000000000" {
				t.Error("span id is zero value")
			}
		}))

		mid := logz.HTTPMiddleware("test/component")(mux)
		req1 := httptest.NewRequest(http.MethodGet, "/test1", nil)

		// Simulates managed cloud service like App Engine or Cloud Run, that sets HTTP header of X-Cloud-Trace-Context
		req1.Header.Set("X-Cloud-Trace-Context", "a0d3eee13de6a4bbcf291eb444b94f28/913411593c9338c5;o=1")

		rec1 := httptest.NewRecorder()
		mid.ServeHTTP(rec1, req1)
	})

}
