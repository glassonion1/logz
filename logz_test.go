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

func TestLogzInfof(t *testing.T) {

	now := time.Date(2020, 12, 31, 23, 59, 59, 999999999, time.Local)
	logz.SetNow(now)

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
		logz.Infof(context.Background(), "writes %s log", "info")

		w.Close()

		var buf bytes.Buffer
		buf.ReadFrom(r)

		// Gets the log from buffer.
		got := strings.TrimRight(buf.String(), "\n")

		expected := `{"severity":"200","message":"writes info log","time":"2020-12-31T23:59:59.999999999+09:00","logging.googleapis.com/trace":"","logging.googleapis.com/spanId":"","jsonPayload":null}`

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
		logz.Warningf(context.Background(), "writes %s log", "warning")

		w.Close()

		var buf bytes.Buffer
		buf.ReadFrom(r)

		// Gets the log from buffer.
		got := strings.TrimRight(buf.String(), "\n")

		expected := `{"severity":"400","message":"writes warning log","time":"2020-12-31T23:59:59.999999999+09:00","logging.googleapis.com/trace":"","logging.googleapis.com/spanId":"","jsonPayload":null}`

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
		logz.Errorf(context.Background(), "writes %s log", "error")

		w.Close()

		var buf bytes.Buffer
		buf.ReadFrom(r)

		// Gets the log from buffer.
		got := strings.TrimRight(buf.String(), "\n")

		expected := `{"severity":"500","message":"writes error log","time":"2020-12-31T23:59:59.999999999+09:00","logging.googleapis.com/trace":"","logging.googleapis.com/spanId":"","jsonPayload":null}`

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
		logz.Criticalf(context.Background(), "writes %s log", "critical")

		w.Close()

		var buf bytes.Buffer
		buf.ReadFrom(r)

		// Gets the log from buffer.
		got := strings.TrimRight(buf.String(), "\n")

		expected := `{"severity":"600","message":"writes critical log","time":"2020-12-31T23:59:59.999999999+09:00","logging.googleapis.com/trace":"","logging.googleapis.com/spanId":"","jsonPayload":null}`

		if diff := cmp.Diff(got, expected); diff != "" {
			// Restores the stdout
			os.Stdout = orgStdout
			t.Errorf("failed test: %v", diff)
		}
	})
}
