package logz

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// HTTPMiddleware is middleware for HTTP handler
func HTTPMiddleware(label string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {

		return otelhttp.NewHandler(h, label)
	}
}
