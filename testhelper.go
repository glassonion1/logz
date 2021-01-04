package logz

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func ExtractStdout(t *testing.T, fnc func()) string {
	// Evacuates the stderr
	orgStdout := os.Stdout
	defer func() {
		os.Stdout = orgStdout
	}()

	// Overrides the stderr to the buffer.
	r, w, _ := os.Pipe()
	os.Stdout = w

	fnc()

	w.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read buf: %v", err)
	}

	return strings.TrimRight(buf.String(), "\n")
}

func ExtractStderr(t *testing.T, fnc func()) string {
	// Evacuates the stderr
	orgStderr := os.Stderr
	defer func() {
		os.Stderr = orgStderr
	}()

	// Overrides the stderr to the buffer.
	r, w, _ := os.Pipe()
	os.Stderr = w

	fnc()

	w.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read buf: %v", err)
	}

	return strings.TrimRight(buf.String(), "\n")
}
