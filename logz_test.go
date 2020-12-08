package logz_test

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/glassonion1/logz"
	"github.com/google/go-cmp/cmp"
)

/*
Tests logz functions.
The log format is below.
{
    "severity":"200",
    "message":"writes info log",
    "time":"2020-12-31T23:59:59.999999999+09:00",
    "logging.googleapis.com/trace":"projects/test/traces/00000000000000000000000000000000",
    "logging.googleapis.com/spanId":"0000000000000000",
    "jsonPayload":null
}
*/
func TestLogz(t *testing.T) {

	ctx := context.Background()

	now := time.Date(2020, 12, 31, 23, 59, 59, 999999999, time.Local)
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

		expected := `{"severity":"200","message":"writes info log","time":"2020-12-31T23:59:59.999999999+09:00","logging.googleapis.com/trace":"projects/test/traces/00000000000000000000000000000000","logging.googleapis.com/spanId":"0000000000000000","jsonPayload":null}`

		if diff := cmp.Diff(got, expected); diff != "" {
			// Restores the stdout
			os.Stdout = orgStdout
			t.Errorf("failed test: %v", diff)
		}
	})

	t.Run("Tests the Warningf function", func(t *testing.T) {
		// Overrides the stdout to the buffer.
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Tests the function
		logz.Warningf(ctx, "writes %s log", "warning")

		w.Close()

		var buf bytes.Buffer
		if _, err := buf.ReadFrom(r); err != nil {
			t.Fatalf("failed to read buf: %v", err)
		}

		// Gets the log from buffer.
		got := strings.TrimRight(buf.String(), "\n")

		expected := `{"severity":"400","message":"writes warning log","time":"2020-12-31T23:59:59.999999999+09:00","logging.googleapis.com/trace":"projects/test/traces/00000000000000000000000000000000","logging.googleapis.com/spanId":"0000000000000000","jsonPayload":null}`

		if diff := cmp.Diff(got, expected); diff != "" {
			// Restores the stdout
			os.Stdout = orgStdout
			t.Errorf("failed test: %v", diff)
		}
	})

	t.Run("Tests the Errorf function", func(t *testing.T) {
		// Overrides the stdout to the buffer.
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Tests the function
		logz.Errorf(ctx, "writes %s log", "error")

		w.Close()

		var buf bytes.Buffer
		if _, err := buf.ReadFrom(r); err != nil {
			t.Fatalf("failed to read buf: %v", err)
		}

		// Gets the log from buffer.
		got := strings.TrimRight(buf.String(), "\n")

		expected := `{"severity":"500","message":"writes error log","time":"2020-12-31T23:59:59.999999999+09:00","logging.googleapis.com/trace":"projects/test/traces/00000000000000000000000000000000","logging.googleapis.com/spanId":"0000000000000000","jsonPayload":null}`

		if diff := cmp.Diff(got, expected); diff != "" {
			// Restores the stdout
			os.Stdout = orgStdout
			t.Errorf("failed test: %v", diff)
		}
	})

	t.Run("Tests the Criticalf function", func(t *testing.T) {
		// Overrides the stdout to the buffer.
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Tests the function
		logz.Criticalf(ctx, "writes %s log", "critical")

		w.Close()

		var buf bytes.Buffer
		if _, err := buf.ReadFrom(r); err != nil {
			t.Fatalf("failed to read buf: %v", err)
		}

		// Gets the log from buffer.
		got := strings.TrimRight(buf.String(), "\n")

		expected := `{"severity":"600","message":"writes critical log","time":"2020-12-31T23:59:59.999999999+09:00","logging.googleapis.com/trace":"projects/test/traces/00000000000000000000000000000000","logging.googleapis.com/spanId":"0000000000000000","jsonPayload":null}`

		if diff := cmp.Diff(got, expected); diff != "" {
			// Restores the stdout
			os.Stdout = orgStdout
			t.Errorf("failed test: %v", diff)
		}
	})
}
