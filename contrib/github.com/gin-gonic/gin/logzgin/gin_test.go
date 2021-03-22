package logzgin_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glassonion1/logz"
	"github.com/glassonion1/logz/contrib/github.com/gin-gonic/gin/logzgin"
	"github.com/glassonion1/logz/testhelper"
	"github.com/google/go-cmp/cmp"
	"go.opentelemetry.io/otel/trace"
)

func TestMiddleware(t *testing.T) {

	logz.InitTracer()

	t.Run("Tests logzgin integration", func(t *testing.T) {
		t.Parallel()

		router := gin.New()
		router.Use(logzgin.Middleware("foobar"))
		router.GET("/test1", func(c *gin.Context) {
			r := c.Request

			ctx := r.Context()
			logz.Infof(ctx, "write %s log", "info")

			sc := trace.SpanContextFromContext(ctx)

			if diff := cmp.Diff(sc.TraceID().String(), "a0d3eee13de6a4bbcf291eb444b94f28"); diff != "" {
				t.Errorf("remote and current trace id are missmatch: %v", diff)
			}
			if sc.SpanID().String() == "0000000000000000" {
				t.Error("span id is zero value")
			}
		})

		r := httptest.NewRequest(http.MethodGet, "/test1", nil)
		// Simulates managed cloud service like App Engine or Cloud Run, that sets HTTP header of X-Cloud-Trace-Context
		r.Header.Set("X-Cloud-Trace-Context", "a0d3eee13de6a4bbcf291eb444b94f28/1;o=1")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
	})
}

func TestMiddlewareMaxSeverity(t *testing.T) {

	t.Run("Tests logzgin integration", func(t *testing.T) {
		t.Parallel()

		got := testhelper.ExtractAccessLogOut(t, func() {
			router := gin.New()
			router.Use(logzgin.Middleware("foobar"))
			router.GET("/test1", func(c *gin.Context) {
				r := c.Request
				ctx := r.Context()

				logz.Infof(ctx, "write %s log", "info")
				logz.Errorf(ctx, "write %s log", "error")
			})
			router.GET("/test2", func(c *gin.Context) {
				r := c.Request
				ctx := r.Context()

				logz.Debugf(ctx, "write %s log1", "debug")
				logz.Debugf(ctx, "write %s log2", "debug")
			})
			router.GET("/test3", func(c *gin.Context) {
				r := c.Request
				ctx := r.Context()

				logz.Warningf(ctx, "write %s log", "warning")
			})

			rec := httptest.NewRecorder()

			r1 := httptest.NewRequest(http.MethodGet, "/test1", nil)
			router.ServeHTTP(rec, r1)
			r2 := httptest.NewRequest(http.MethodGet, "/test2", nil)
			router.ServeHTTP(rec, r2)
			r3 := httptest.NewRequest(http.MethodGet, "/test3", nil)
			router.ServeHTTP(rec, r3)
		})

		if !strings.Contains(got, `"severity":"ERROR"`) {
			t.Error("max severity is not set correctly")
		}
		if !strings.Contains(got, `"severity":"INFO"`) {
			t.Error("max severity is not set correctly")
		}
		if !strings.Contains(got, `"severity":"WARNING"`) {
			t.Error("max severity is not set correctly")
		}
	})
}
