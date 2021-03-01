package go111

import (
	"context"
	"io"

	"github.com/glassonion1/logz/go111/internal/config"
	"github.com/glassonion1/logz/go111/internal/logger"
	"github.com/glassonion1/logz/go111/internal/severity"
)

// Config is configurations for logz
type Config struct {
	// GCP Project ID
	ProjectID string
	// CallerDepth is the number of stack frames to ascend
	CallerSkip int
	// Output for application log
	ApplicationLogOut io.Writer
}

// SetProjectID sets gcp project id to the logger
func SetProjectID(projectID string) {
	config.ProjectID = projectID
}

// SetConfig sets config to the logger
func SetConfig(cfg Config) {
	if cfg.ProjectID != "" {
		config.ProjectID = cfg.ProjectID
	}

	config.CallerSkip = cfg.CallerSkip

	if cfg.ApplicationLogOut != nil {
		config.ApplicationLogOut = cfg.ApplicationLogOut
	}
}

// Debugf writes debug log to the stdout
func Debugf(ctx context.Context, format string, a ...interface{}) {
	logger.WriteApplicationLog(ctx, severity.Debug, format, a...)
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
