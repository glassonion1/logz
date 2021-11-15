package testhelper

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/glassonion1/logz/internal/config"
)

// ExtractStdout extracts string from stdout
func ExtractApplicationLogOut(t *testing.T, ctx context.Context, fnc func(ctx context.Context)) string {
	t.Helper()

	var buf bytes.Buffer
	if cc := config.GetContextConfig(ctx); cc != nil {
		cc.ApplicationLogOut = &buf
	} else {
		ctx = config.SetContextConfig(ctx, &config.ContextConfig{ApplicationLogOut: &buf})
	}
	fnc(ctx)

	return strings.TrimRight(buf.String(), "\n")
}

// ExtractStdout extracts string from stderr
func ExtractAccessLogOut(t *testing.T, ctx context.Context, fnc func(ctx context.Context)) string {
	t.Helper()

	var buf bytes.Buffer
	if cc := config.GetContextConfig(ctx); cc != nil {
		cc.AccessLogOut = &buf
	} else {
		ctx = config.SetContextConfig(ctx, &config.ContextConfig{AccessLogOut: &buf})
	}

	fnc(ctx)

	return strings.TrimRight(buf.String(), "\n")
}

// OverrideLogOutContext override log I/O in the context
func OverrideLogOutContext(t *testing.T, ctx context.Context, appLogOut, accessLogOut io.Writer) context.Context {
	t.Helper()
	return config.SetContextConfig(ctx, &config.ContextConfig{
		ApplicationLogOut: appLogOut,
		AccessLogOut:      accessLogOut,
	})
}
