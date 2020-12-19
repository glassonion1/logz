package logz

import (
	"context"
	"net/http"
	"time"

	"github.com/glassonion1/logz/internal/config"
	"github.com/glassonion1/logz/internal/logger"
	"github.com/glassonion1/logz/internal/severity"
	logzpropagation "github.com/glassonion1/logz/propagation"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// SetProjectID sets gcp project id to the logger
func SetProjectID(projectID string) {
	config.ProjectID = projectID
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

// Access writes access log to the stderr
func Access(ctx context.Context, r http.Request, statusCode, size int, elapsed time.Duration) {
	logger.WriteAccessLog(ctx, r, statusCode, size, elapsed)
}

// InitTracer initializes OpenTelemetry tracer
func InitTracer() {
	tp := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(tp)

	props := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
		logzpropagation.HTTPFormat{})
	otel.SetTextMapPropagator(props)
}

// StartCollectingSeverity starts collectiong severity
func StartCollectingSeverity(ctx context.Context) context.Context {
	cs := &severity.ContextSeverity{}
	return severity.SetContextSeverity(ctx, cs)
}
