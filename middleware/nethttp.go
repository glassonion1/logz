package middleware

import (
	"net/http"
	"time"

	"github.com/glassonion1/logz/internal/logger"
	"github.com/glassonion1/logz/writer"
	"go.opentelemetry.io/otel"
)

// NetHTTP is middleware for HTTP handler
func NetHTTP(label string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := logger.NowFunc()
			rw := writer.NewResponseWriter(w)

			tracer := otel.Tracer(label)

			prop := otel.GetTextMapPropagator()
			ctx := prop.Extract(r.Context(), r.Header)

			newCtx, span := tracer.Start(ctx, r.URL.String())

			defer func() {
				tID := span.SpanContext().TraceID.String()
				logger.WriteAccessLog(tID, *r, rw.StatusCode(), rw.Size(), time.Since(started))
				span.End()
			}()

			h.ServeHTTP(rw, r.WithContext(newCtx))
		})
	}
}
