package logz

import (
	"net/http"

	"github.com/glassonion1/logz/internal/tracer"
)

// HTTPMiddleware is middleware for HTTP handler
func HTTPMiddleware(label string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tra := tracer.New(label)
			childCtx, end := tra.Start(r.Context(), r.URL.String())
			defer end()

			h.ServeHTTP(w, r.WithContext(childCtx))
		})
	}
}
