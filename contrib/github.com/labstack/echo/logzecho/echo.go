package logzecho

import (
	"time"

	"github.com/glassonion1/logz"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// Middleware is middleware for HTTP handler
func Middleware(label string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			started := time.Now()

			tracer := otel.Tracer(label)

			r := c.Request()

			prop := otel.GetTextMapPropagator()
			ctx := prop.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			ctx, span := tracer.Start(ctx, r.URL.String())
			ctx = logz.StartCollectingSeverity(ctx)

			defer func() {
				size := int(c.Response().Size)
				logz.Access(ctx, *r, c.Response().Status, size, time.Since(started))
				span.End()
			}()

			// pass the span through the request context
			c.SetRequest(r.WithContext(ctx))

			// serve the request to the next middleware
			return next(c)
		}
	}
}
