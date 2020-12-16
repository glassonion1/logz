package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glassonion1/logz/internal/spancontext"
	"github.com/glassonion1/logz/middleware"
)

func TestNetHTTP(t *testing.T) {

	t.Run("Tests the middleware", func(t *testing.T) {

		mux := http.NewServeMux()
		mux.Handle("/test1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sc := spancontext.Extract(r.Context())
			if sc.TraceID == "" || sc.SpanID == "" {
				t.Error("failed to test middleware, span context is zero value.")
			}
			fmt.Fprintf(w, "hello world")
		}))

		mid := middleware.NetHTTP("test/component")(mux)
		req1 := httptest.NewRequest(http.MethodGet, "/test1", nil)
		rec1 := httptest.NewRecorder()
		mid.ServeHTTP(rec1, req1)
	})

}
