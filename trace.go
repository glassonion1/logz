package logz

import (
	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/api/option"

	"github.com/glassonion1/logz/internal/config"
	logzpropagation "github.com/glassonion1/logz/propagation"
)

func InitTracer() error {
	tp := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(logzpropagation.HTTPFormat{}))
	return nil
}

// InitStdoutTracer initializes tracer of OpenTelemetry, that is for stdout
func InitStdoutTracer() error {
	// Create stdout exporter to be able to retrieve
	// the collected spans.
	exporter, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		return err
	}

	config := sdktrace.Config{
		DefaultSampler: sdktrace.AlwaysSample(),
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(config),
		sdktrace.WithSyncer(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}, logzpropagation.HTTPFormat{}))

	return nil
}

// InitCloudTracer initializes tracer of OpenTelemetry, that is for Cloud Logging(formerly known as Stackdriver Logging)
func InitCloudTracer(opts ...option.ClientOption) error {

	traceOpts := []cloudtrace.Option{
		cloudtrace.WithTraceClientOptions(opts),
		cloudtrace.WithProjectID(config.ProjectID),
	}

	// Create cloud tracer exporter to be able to retrieve
	// the collected spans.
	exporter, err := cloudtrace.NewExporter(traceOpts...)
	if err != nil {
		return err
	}

	config := sdktrace.Config{
		DefaultSampler: sdktrace.AlwaysSample(),
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(config),
		sdktrace.WithSyncer(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}, logzpropagation.HTTPFormat{}))

	return nil
}
