package logger_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/glassonion1/logz/internal/config"
	"github.com/glassonion1/logz/internal/logger"
	"github.com/glassonion1/logz/internal/severity"
	"github.com/glassonion1/logz/internal/types"
	"github.com/glassonion1/logz/testhelper"
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
    "function":"github.com/glassonion1/logz.ExtractApplicationLogOut"
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
	config.CallerSkip = 1

	defer func() {
		config.CallerSkip = 0
	}()

	t.Run("Tests WriteApplicationLog function", func(t *testing.T) {
		got := testhelper.ExtractApplicationLogOut(t, ctx, func(ctx context.Context) {
			// tests the function
			logger.WriteApplicationLog(ctx, severity.Info, "writes %s log", "info")
		})

		expected := `{"severity":"INFO","message":"writes info log","time":"2020-12-31T23:59:59.999999999Z","logging.googleapis.com/sourceLocation":{"file":"logger_test.go","line":"51","function":"github.com/glassonion1/logz/internal/logger_test.TestLoggerWriteApplicationLog.func3"},"logging.googleapis.com/trace":"projects/test/traces/00000000000000000000000000000000","logging.googleapis.com/spanId":"0000000000000000","logging.googleapis.com/trace_sampled":false}`

		if diff := cmp.Diff(got, expected); diff != "" {
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

	ctx := context.Background()

	types.GetServerIP = func() string {
		return "192.168.0.1"
	}

	logger.NowFunc = func() time.Time {
		return time.Date(2020, 12, 31, 23, 59, 59, 999999999, time.UTC)
	}
	config.ProjectID = "test"

	t.Run("Tests WriteAccessLog function", func(t *testing.T) {

		got := testhelper.ExtractAccessLogOut(t, ctx, func(ctx context.Context) {
			// Tests the function
			httpReq := httptest.NewRequest(http.MethodGet, "/test1", nil)
			req := types.MakeHTTPRequest(*httpReq, 200, 333, time.Duration(100))
			logger.WriteAccessLog(ctx, req)
		})

		expected := `{"severity":"INFO","time":"2020-12-31T23:59:59.999999999Z","logging.googleapis.com/trace":"projects/test/traces/00000000000000000000000000000000","httpRequest":{"requestMethod":"GET","requestUrl":"/test1","requestSize":"0","status":200,"responseSize":"333","remoteIp":"192.0.2.1","serverIp":"192.168.0.1","latency":{"nanos":100,"seconds":0},"protocol":"HTTP/1.1"}}`

		if diff := cmp.Diff(got, expected); diff != "" {
			t.Errorf("failed log info test: %v", diff)
		}
	})

}
