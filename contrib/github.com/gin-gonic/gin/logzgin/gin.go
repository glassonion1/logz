package logzgin

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glassonion1/logz"
	"go.opentelemetry.io/otel"
)

// Middleware is middleware for HTTP handler
func Middleware(label string) gin.HandlerFunc {
	return func(c *gin.Context) {

		started := time.Now()

		tracer := otel.Tracer(label)

		r := c.Request

		prop := otel.GetTextMapPropagator()
		ctx := prop.Extract(r.Context(), r.Header)

		ctx, span := tracer.Start(ctx, r.URL.String())
		ctx = logz.StartCollectingSeverity(ctx)

		defer func() {
			logz.Access(ctx, *r, c.Writer.Status(), c.Writer.Size(), time.Since(started))
			span.End()
		}()

		// pass the span through the request context
		c.Request = c.Request.WithContext(ctx)

		// serve the request to the next middleware
		c.Next()
	}
}
