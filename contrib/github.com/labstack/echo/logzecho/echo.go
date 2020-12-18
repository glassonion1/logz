package logzecho

import (
	"time"

	"github.com/glassonion1/logz"
	"github.com/labstack/echo"
	"go.opentelemetry.io/otel"
)

func Middleware(label string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			started := time.Now()

			tracer := otel.Tracer(label)

			r := c.Request()

			prop := otel.GetTextMapPropagator()
			ctx := prop.Extract(r.Context(), r.Header)

			newCtx, span := tracer.Start(ctx, r.URL.String())

			defer func() {
				tID := span.SpanContext().TraceID.String()
				size := int(c.Response().Size)
				logz.Access(tID, *r, c.Response().Status, size, time.Since(started))
				span.End()
			}()

			// pass the span through the request context
			c.SetRequest(r.WithContext(newCtx))

			// serve the request to the next middleware
			return next(c)
		}
	}
}
