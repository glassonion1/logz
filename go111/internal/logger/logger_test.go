package logger_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/glassonion1/logz/go111/internal/config"
	"github.com/glassonion1/logz/go111/internal/logger"
	"github.com/glassonion1/logz/go111/internal/severity"
	"github.com/glassonion1/logz/go111/testhelper"
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
    "function":"github.com/glassonion1/logz/internal/logger_test.TestLoggerWriteApplicationLog.func2"
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

	t.Run("Tests WriteApplicationLog function", func(t *testing.T) {
		got := testhelper.ExtractApplicationLogOut(t, func() {
			// tests the function
			logger.WriteApplicationLog(ctx, severity.Info, "writes %s log", "info")
		})

		expected := `{"severity":"INFO","message":"writes info log","time":"2020-12-31T23:59:59.999999999Z","logging.googleapis.com/sourceLocation":{"file":"logger_test.go","line":"44","function":"github.com/glassonion1/logz/go111/internal/logger_test.TestLoggerWriteApplicationLog.func2"},"logging.googleapis.com/trace":"projects/test/traces/00000000000000000000000000000000","logging.googleapis.com/spanId":"0000000000000000","logging.googleapis.com/trace_sampled":false}`

		if diff := cmp.Diff(got, expected); diff != "" {
			t.Errorf("failed log info test: %v", diff)
		}
	})
}
