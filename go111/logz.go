package go111

import (
	"context"

	"github.com/glassonion1/logz/go111/internal/config"
	"github.com/glassonion1/logz/go111/internal/logger"
	"github.com/glassonion1/logz/go111/internal/severity"
)

// SetProjectID sets gcp project id to the logger
func SetProjectID(projectID string) {
	config.ProjectID = projectID
}

// Debugf writes debug log to the stdout
func Debugf(ctx context.Context, format string, a ...interface{}) {
	logger.WriteApplicationLog(ctx, severity.Default, format, a...)
}

// Infof writes info log to the stdout
func Infof(ctx context.Context, format string, a ...interface{}) {
	logger.WriteApplicationLog(ctx, severity.Info, format, a...)
}

// Warningf writes warning log to the stdout
func Warningf(ctx context.Context, format string, a ...interface{}) {
	logger.WriteApplicationLog(ctx, severity.Warning, format, a...)
}

// Errorf writes error log to the stdout
func Errorf(ctx context.Context, format string, a ...interface{}) {
	logger.WriteApplicationLog(ctx, severity.Error, format, a...)
}

// Criticalf writes critical log to the stdout
func Criticalf(ctx context.Context, format string, a ...interface{}) {
	logger.WriteApplicationLog(ctx, severity.Critical, format, a...)
}
