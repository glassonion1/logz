package logz

import (
	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

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

	return nil
}

func InitCloudTracer(projectID string, opts ...cloudtrace.Option) error {

	opts = append(opts, cloudtrace.WithProjectID(projectID))

	// Create cloud tracer exporter to be able to retrieve
	// the collected spans.
	exporter, err := cloudtrace.NewExporter(opts...)
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

	return nil
}
