package logger_test

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"context"

	"github.com/glassonion1/logz/internal/config"
	"github.com/glassonion1/logz/internal/logger"
	"github.com/glassonion1/logz/internal/severity"
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
    "file":"testing.go",
    "line":"1127",
    "function":"testing.tRunner"
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
		logger.WriteApplicationLog(ctx, severity.Info, "writes %s log", "info")

		w.Close()

		var buf bytes.Buffer
		if _, err := buf.ReadFrom(r); err != nil {
			t.Fatalf("failed to read buf: %v", err)
		}

		// Gets the log from buffer.
		got := strings.TrimRight(buf.String(), "\n")

		expected := `{"severity":"INFO","message":"writes info log","time":"2020-12-31T23:59:59.999999999Z","logging.googleapis.com/sourceLocation":{"file":"testing.go","line":"1127","function":"testing.tRunner"},"logging.googleapis.com/trace":"projects/test/traces/00000000000000000000000000000000","logging.googleapis.com/spanId":"0000000000000000","logging.googleapis.com/trace_sampled":false}`

		if diff := cmp.Diff(got, expected); diff != "" {
			// Restores the stdout
			os.Stdout = orgStdout
			t.Errorf("failed log info test: %v", diff)
		}
	})
}
