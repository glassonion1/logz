package testhelper

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/glassonion1/logz/internal/config"
)

// ExtractStdout extracts string from stdout
func ExtractApplicationLogOut(t *testing.T, fnc func()) string {
	t.Helper()

	var buf bytes.Buffer
	config.ApplicationLogOut = &buf
	defer func() {
		config.ApplicationLogOut = os.Stdout
	}()

	fnc()

	return strings.TrimRight(buf.String(), "\n")
}

// ExtractStdout extracts string from stderr
func ExtractAccessLogOut(t *testing.T, fnc func()) string {
	t.Helper()

	var buf bytes.Buffer
	config.AccessLogOut = &buf
	defer func() {
		config.AccessLogOut = os.Stdout
	}()

	fnc()

	return strings.TrimRight(buf.String(), "\n")
}
