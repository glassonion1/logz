package logz

import (
	"context"

	"github.com/glassonion1/logz/internal/logger"
	"github.com/glassonion1/logz/internal/severity"
)

var std = &logger.Logger{}

func Debugf(ctx context.Context, format string, a ...interface{}) {
	std.WriteLog(ctx, severity.Default, format, a...)
}

func Infof(ctx context.Context, format string, a ...interface{}) {
	std.WriteLog(ctx, severity.Info, format, a...)
}

func Warningf(ctx context.Context, format string, a ...interface{}) {
	std.WriteLog(ctx, severity.Warning, format, a...)
}

func Errorf(ctx context.Context, format string, a ...interface{}) {
	std.WriteLog(ctx, severity.Error, format, a...)
}

func Criticalf(ctx context.Context, format string, a ...interface{}) {
	std.WriteLog(ctx, severity.Critical, format, a...)
}
