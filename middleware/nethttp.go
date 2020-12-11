package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel"
)

// NetHTTP is middleware for HTTP handler
func NetHTTP(label string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tracer := otel.Tracer(label)

			prop := otel.GetTextMapPropagator()
			ctx := prop.Extract(r.Context(), r.Header)

			newCtx, span := tracer.Start(ctx, r.URL.String())
			defer span.End()

			h.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}
