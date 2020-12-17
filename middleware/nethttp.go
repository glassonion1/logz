package middleware

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel"

	"github.com/glassonion1/logz"
	"github.com/glassonion1/logz/writer"
)

// NetHTTP is middleware for HTTP handler
func NetHTTP(label string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			rw := writer.NewResponseWriter(w)

			tracer := otel.Tracer(label)

			prop := otel.GetTextMapPropagator()
			ctx := prop.Extract(r.Context(), r.Header)

			newCtx, span := tracer.Start(ctx, r.URL.String())

			defer func() {
				tID := span.SpanContext().TraceID.String()
				logz.Access(tID, *r, rw.StatusCode(), rw.Size(), time.Since(started))
				span.End()
			}()

			h.ServeHTTP(rw, r.WithContext(newCtx))
		})
	}
}
