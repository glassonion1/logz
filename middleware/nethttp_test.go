package middleware_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/glassonion1/logz"
	"github.com/glassonion1/logz/internal/spancontext"
	"github.com/glassonion1/logz/middleware"
)

func TestNetHTTP(t *testing.T) {

	t.Run("Tests the middleware", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.Handle("/test1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			sc := spancontext.Extract(ctx)
			if sc.TraceID == "" || sc.SpanID == "" {
				t.Error("failed to test middleware, span context is zero value.")
			}

			fmt.Fprintf(w, "hello world")
		}))

		mid := middleware.NetHTTP("test/component")(mux)
		req := httptest.NewRequest(http.MethodGet, "/test1", nil)
		rec := httptest.NewRecorder()
		mid.ServeHTTP(rec, req)
	})
}

func TestNetHTTPMaxSeverity(t *testing.T) {

	// Evacuates the stderr
	orgStderr := os.Stderr
	defer func() {
		os.Stderr = orgStderr
	}()
	t.Run("Tests the middleware", func(t *testing.T) {
		// Overrides the stderr to the buffer.
		r, w, _ := os.Pipe()
		os.Stderr = w

		mux := http.NewServeMux()
		mux.Handle("/test1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			logz.Infof(ctx, "write %s log", "info")
			logz.Errorf(ctx, "write %s log", "error")

			fmt.Fprintf(w, "hello world")
		}))
		mux.Handle("/test2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			logz.Debugf(ctx, "write %s log1", "debug")
			logz.Debugf(ctx, "write %s log2", "debug")

			fmt.Fprintf(w, "hello world")
		}))
		mux.Handle("/test3", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			logz.Warningf(ctx, "write %s log", "warning")

			fmt.Fprintf(w, "hello world")
		}))

		mid := middleware.NetHTTP("test/component")(mux)
		rec := httptest.NewRecorder()

		req1 := httptest.NewRequest(http.MethodGet, "/test1", nil)
		mid.ServeHTTP(rec, req1)
		req2 := httptest.NewRequest(http.MethodGet, "/test2", nil)
		mid.ServeHTTP(rec, req2)
		req3 := httptest.NewRequest(http.MethodGet, "/test3", nil)
		mid.ServeHTTP(rec, req3)

		// Tests max severity of access log
		w.Close()
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(r); err != nil {
			t.Fatalf("failed to read buf: %v", err)
		}

		// Gets the log from buffer.
		got := buf.String()
		if !strings.Contains(got, `"severity":"ERROR"`) {
			t.Error("max severity is not set correctly: error")
		}
		if !strings.Contains(got, `"severity":"INFO"`) {
			t.Error("max severity is not set correctly: info")
		}
		if !strings.Contains(got, `"severity":"WARNING"`) {
			t.Error("max severity is not set correctly: warning")
		}
	})
}

func TestNetHTTPMaxSeverityNoLog(t *testing.T) {

	// Evacuates the stderr
	orgStderr := os.Stderr
	defer func() {
		os.Stderr = orgStderr
	}()
	t.Run("Tests the middleware", func(t *testing.T) {
		// Overrides the stderr to the buffer.
		r, w, _ := os.Pipe()
		os.Stderr = w

		mux := http.NewServeMux()
		mux.Handle("/test1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "hello world")
		}))
		mux.Handle("/test2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "hello world")
		}))

		mid := middleware.NetHTTP("test/component")(mux)
		rec := httptest.NewRecorder()

		req1 := httptest.NewRequest(http.MethodGet, "/test1", nil)
		mid.ServeHTTP(rec, req1)
		req2 := httptest.NewRequest(http.MethodGet, "/test2", nil)
		mid.ServeHTTP(rec, req2)

		// Tests max severity of access log
		w.Close()
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(r); err != nil {
			t.Fatalf("failed to read buf: %v", err)
		}

		// Gets the log from buffer.
		got := buf.String()
		if !strings.Contains(got, `"severity":"ERROR"`) {
			t.Error("max severity is not set correctly: error")
		}
		if !strings.Contains(got, `"severity":"INFO"`) {
			t.Error("max severity is not set correctly: info")
		}
	})
}
