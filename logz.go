/*
Package logz provides the structured log.
Example:
	ctx := r.Context() // r is *http.Request
	logz.Infof(ctx, "info log. requestURL: %s", r.URL.String())
*/
package logz

import (
	"context"

	"github.com/glassonion1/logz/internal/logger"
	"github.com/glassonion1/logz/internal/severity"
)

var std = &logger.Logger{}

// Debugf writes debug log to the stdout
func Debugf(ctx context.Context, format string, a ...interface{}) {
	std.WriteLog(ctx, severity.Default, format, a...)
}

// Infof writes info log to the stdout
func Infof(ctx context.Context, format string, a ...interface{}) {
	std.WriteLog(ctx, severity.Info, format, a...)
}

// Warningf writes warning log to the stdout
func Warningf(ctx context.Context, format string, a ...interface{}) {
	std.WriteLog(ctx, severity.Warning, format, a...)
}

// Errorf writes error log to the stdout
func Errorf(ctx context.Context, format string, a ...interface{}) {
	std.WriteLog(ctx, severity.Error, format, a...)
}

// Criticalf writes critical log to the stdout
func Criticalf(ctx context.Context, format string, a ...interface{}) {
	std.WriteLog(ctx, severity.Critical, format, a...)
}
