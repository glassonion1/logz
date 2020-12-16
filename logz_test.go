package logz_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"context"

	"github.com/glassonion1/logz"
	"github.com/glassonion1/logz/internal/spancontext"
	"github.com/glassonion1/logz/middleware"
	"github.com/google/go-cmp/cmp"
)

func TestLogz(t *testing.T) {

	logz.InitTracer()

	t.Run("Tests logz integration", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		mux.Handle("/test1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logz.Infof(ctx, "write %s log", "info")

			sc := spancontext.Extract(ctx)

			if sc.TraceID == "00000000000000000000000000000000" {
				t.Error("trace id is zero value")
			}
			if sc.SpanID == "0000000000000000" {
				t.Error("span id is zero value")
			}
		}))

		mux.Handle("/test2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logz.Infof(ctx, "write %s log", "info")

			sc := spancontext.Extract(ctx)

			if sc.TraceID == "00000000000000000000000000000000" {
				t.Error("trace id is zero value")
			}
			if sc.SpanID == "0000000000000000" {
				t.Error("span id is zero value")
			}

			// nested function
			func(ctx context.Context) {
				logz.Infof(ctx, "write %s nested log", "info")

				child := spancontext.Extract(ctx)
				if child.TraceID != sc.TraceID {
					t.Error("trace and child trace id are not equal")
				}
				if child.SpanID != sc.SpanID {
					t.Error("span and child span id are not equal")
				}
			}(ctx)
		}))

		mid := middleware.NetHTTP("test/component")(mux)
		req1 := httptest.NewRequest(http.MethodGet, "/test1", nil)
		rec1 := httptest.NewRecorder()
		mid.ServeHTTP(rec1, req1)

		req2 := httptest.NewRequest(http.MethodGet, "/test2", nil)
		rec2 := httptest.NewRecorder()
		mid.ServeHTTP(rec2, req2)
	})

}

func TestLogzRemoteParent(t *testing.T) {

	logz.InitTracer()

	t.Run("Tests logz integration with remote parent", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()
		mux.Handle("/test1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logz.Infof(ctx, "write %s log", "info")

			sc := spancontext.Extract(ctx)

			if diff := cmp.Diff(sc.TraceID, "a0d3eee13de6a4bbcf291eb444b94f28"); diff != "" {
				t.Errorf("remote and current trace id are missmatch: %v", diff)
			}
			if sc.SpanID == "0000000000000000" {
				t.Error("span id is zero value")
			}
		}))

		mid := middleware.NetHTTP("test/component")(mux)
		req1 := httptest.NewRequest(http.MethodGet, "/test1", nil)

		// Simulates managed cloud service like App Engine or Cloud Run, that sets HTTP header of X-Cloud-Trace-Context
		req1.Header.Set("X-Cloud-Trace-Context", "a0d3eee13de6a4bbcf291eb444b94f28/1;o=1")

		rec1 := httptest.NewRecorder()
		mid.ServeHTTP(rec1, req1)
	})

}
