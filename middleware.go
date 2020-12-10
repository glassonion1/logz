package logz

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

// HTTPMiddleware is middleware for HTTP handler
func HTTPMiddleware2(label string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {

		return otelhttp.NewHandler(h, label)
	}
}

func HTTPMiddleware(label string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tracer := otel.Tracer("logz")

			prop := otel.GetTextMapPropagator()
			ctx := prop.Extract(r.Context(), r.Header)

			newCtx, span := tracer.Start(ctx, r.URL.String())
			defer span.End()

			h.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}
