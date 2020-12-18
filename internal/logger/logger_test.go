package logger_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"context"

	"github.com/glassonion1/logz/internal/config"
	"github.com/glassonion1/logz/internal/logger"
	"github.com/glassonion1/logz/internal/severity"
	"github.com/glassonion1/logz/internal/types"
	"github.com/google/go-cmp/cmp"
)

/*
Tests logger functions.
The format of application log is below
{
  "severity":"INFO",
  "message":"writes info log",
  "time":"2020-12-31T23:59:59.999999999Z",
  "logging.googleapis.com/sourceLocation":{
    "file":"logger_test.go",
    "line":"57",
    "function":"github.com/glassonion1/logz/internal/logger_test.TestLoggerWriteApplicationLog.func3"
  },
  "logging.googleapis.com/trace":"projects/test/traces/00000000000000000000000000000000",
  "logging.googleapis.com/spanId":"0000000000000000",
  "logging.googleapis.com/trace_sampled":false
}
*/
func TestLoggerWriteApplicationLog(t *testing.T) {

	ctx := context.Background()

	logger.NowFunc = func() time.Time {
		return time.Date(2020, 12, 31, 23, 59, 59, 999999999, time.UTC)
	}
	config.ProjectID = "test"

	// Evacuates the stdout
	orgStdout := os.Stdout
	defer func() {
		os.Stdout = orgStdout
	}()
	t.Run("Tests WriteApplicationLog function", func(t *testing.T) {
		// Overrides the stdout to the buffer.
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Tests the function
		func() {
			// Wrapped in a function to have the source location output the intended string
			logger.WriteApplicationLog(ctx, severity.Info, "writes %s log", "info")
		}()

		w.Close()

		var buf bytes.Buffer
		if _, err := buf.ReadFrom(r); err != nil {
			t.Fatalf("failed to read buf: %v", err)
		}

		// Gets the log from buffer.
		got := strings.TrimRight(buf.String(), "\n")

		expected := `{"severity":"INFO","message":"writes info log","time":"2020-12-31T23:59:59.999999999Z","logging.googleapis.com/sourceLocation":{"file":"logger_test.go","line":"61","function":"github.com/glassonion1/logz/internal/logger_test.TestLoggerWriteApplicationLog.func3"},"logging.googleapis.com/trace":"projects/test/traces/00000000000000000000000000000000","logging.googleapis.com/spanId":"0000000000000000","logging.googleapis.com/trace_sampled":false}`

		if diff := cmp.Diff(got, expected); diff != "" {
			// Restores the stdout
			os.Stdout = orgStdout
			t.Errorf("failed log info test: %v", diff)
		}
	})
}

/*
Tests logger functions.
The format of access log is below
{
  "severity":"DEFAULT",
  "time":"2020-12-31T23:59:59.999999999Z",
  "logging.googleapis.com/trace":"projects/test/traces/a0d3eee13de6a4bbcf291eb444b94f28",
  "httpRequest":{
    "requestMethod":"GET",
    "requestUrl":"/test1",
    "requestSize":"0",
    "status":200,
    "responseSize":"333",
    "remoteIp":"192.0.2.1",
    "serverIp":"192.168.100.115",
    "latencyy":{
      "nanos":100,
      "seconds":0
    },
    "protocol":"HTTP/1.1"
  }
}
*/
func TestLoggerWriteAccessLog(t *testing.T) {

	types.GetServerIP = func() string {
		return "192.168.0.1"
	}

	logger.NowFunc = func() time.Time {
		return time.Date(2020, 12, 31, 23, 59, 59, 999999999, time.UTC)
	}
	config.ProjectID = "test"

	// Evacuates the stderr
	orgStderr := os.Stderr
	defer func() {
		os.Stderr = orgStderr
	}()
	t.Run("Tests WriteAccessLog function", func(t *testing.T) {
		// Overrides the stderr to the buffer.
		r, w, _ := os.Pipe()
		os.Stderr = w

		// Tests the function
		req := httptest.NewRequest(http.MethodGet, "/test1", nil)
		logger.WriteAccessLog("a0d3eee13de6a4bbcf291eb444b94f28", *req, 200, 333, time.Duration(100))

		w.Close()

		var buf bytes.Buffer
		if _, err := buf.ReadFrom(r); err != nil {
			t.Fatalf("failed to read buf: %v", err)
		}

		// Gets the log from buffer.
		got := strings.TrimRight(buf.String(), "\n")

		expected := `{"severity":"DEFAULT","time":"2020-12-31T23:59:59.999999999Z","logging.googleapis.com/trace":"projects/test/traces/a0d3eee13de6a4bbcf291eb444b94f28","httpRequest":{"requestMethod":"GET","requestUrl":"/test1","requestSize":"0","status":200,"responseSize":"333","remoteIp":"192.0.2.1","serverIp":"192.168.0.1","latencyy":{"nanos":100,"seconds":0},"protocol":"HTTP/1.1"}}`

		if diff := cmp.Diff(got, expected); diff != "" {
			// Restores the stderr
			os.Stderr = orgStderr
			t.Errorf("failed log info test: %v", diff)
		}
	})
}
