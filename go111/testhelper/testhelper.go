package testhelper

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/glassonion1/logz/go111/internal/config"
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
