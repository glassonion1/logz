package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glassonion1/logz/go111"
	"github.com/glassonion1/logz/go111/internal/spancontext"
	"github.com/glassonion1/logz/go111/middleware"
	"github.com/google/go-cmp/cmp"
)

func TestHTTPMiddleware(t *testing.T) {

	t.Run("Tests the middleware with remote parent", func(t *testing.T) {
		defer func() {

		}()

		mux := http.NewServeMux()
		mux.Handle("/test1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			go111.Infof(ctx, "write %s log", "info")

			sc := spancontext.Extract(ctx)

			if diff := cmp.Diff(sc.TraceID, "a0d3eee13de6a4bbcf291eb444b94f28"); diff != "" {
				t.Errorf("remote and current trace id are missmatch: %v", diff)
			}
			if sc.SpanID == "0000000000000000" {
				t.Errorf("span id is zero value")
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
