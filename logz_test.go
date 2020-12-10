package logz_test

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/glassonion1/logz"
)

/*
Tests logz functions.
The log format is below.
{
    "severity":"INFO",
    "message":"writes info log",
    "time":"2020-12-31T23:59:59.999999999+09:00",
    "logging.googleapis.com/trace":"projects/test/traces/00000000000000000000000000000000",
    "logging.googleapis.com/spanId":"0000000000000000",
    "insertId":"41d8c99e-3ac9-11eb-938c-acde48001122"
}
*/
func TestLogz(t *testing.T) {

	ctx := context.Background()

	now := time.Date(2020, 12, 31, 23, 59, 59, 999999999, time.UTC)
	logz.SetNow(now)
	logz.SetProjectID("test")

	// Evacuates the stdout
	orgStdout := os.Stdout
	defer func() {
		os.Stdout = orgStdout
	}()
	t.Run("Tests the Infof function", func(t *testing.T) {
		// Overrides the stdout to the buffer.
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Tests the function
		logz.Infof(ctx, "writes %s log", "info")

		w.Close()

		var buf bytes.Buffer
		if _, err := buf.ReadFrom(r); err != nil {
			t.Fatalf("failed to read buf: %v", err)
		}

		// Gets the log from buffer.
		got := strings.TrimRight(buf.String(), "\n")

		expected := `"severity":"INFO","message":"writes info log","time":"2020-12-31T23:59:59.999999999Z","logging.googleapis.com/trace":"projects/test/traces/00000000000000000000000000000000","logging.googleapis.com/spanId":"0000000000000000"`

		if !strings.Contains(got, expected) {
			// Restores the stdout
			os.Stdout = orgStdout
			t.Error("failed log info test")
		}

		if !strings.Contains(got, `"insertId":`) {
			// Restores the stdout
			os.Stdout = orgStdout
			t.Error("failed log info test")
		}
	})
}
