package logzecho_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glassonion1/logz"
	"github.com/glassonion1/logz/contrib/github.com/labstack/echo/logzecho"
	"github.com/google/go-cmp/cmp"
	"github.com/labstack/echo"
	"go.opentelemetry.io/otel/trace"
)

func TestMiddleware(t *testing.T) {

	logz.InitTracer()

	t.Run("Tests logzgin integration", func(t *testing.T) {
		t.Parallel()

		router := echo.New()
		router.Use(logzecho.Middleware("foobar"))
		router.GET("/test1", func(c echo.Context) error {
			r := c.Request()

			ctx := r.Context()
			logz.Infof(ctx, "write %s log", "info")

			sc := trace.SpanContextFromContext(ctx)

			if diff := cmp.Diff(sc.TraceID.String(), "a0d3eee13de6a4bbcf291eb444b94f28"); diff != "" {
				t.Errorf("remote and current trace id are missmatch: %v", diff)
			}
			if sc.SpanID.String() == "0000000000000000" {
				t.Error("span id is zero value")
			}
			return c.NoContent(200)
		})

		r := httptest.NewRequest(http.MethodGet, "/test1", nil)
		// Simulates managed cloud service like App Engine or Cloud Run, that sets HTTP header of X-Cloud-Trace-Context
		r.Header.Set("X-Cloud-Trace-Context", "a0d3eee13de6a4bbcf291eb444b94f28/1;o=1")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
	})
}
