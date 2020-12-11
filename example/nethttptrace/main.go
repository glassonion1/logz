package main

import (
	"fmt"
	"log"
	"net/http"

	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/glassonion1/logz"
	logzpropagation "github.com/glassonion1/logz/propagation"
	"github.com/googleinterns/cloud-operations-api-mock/cloudmock"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/api/option"
)

// initCloudTracer initializes tracer of OpenTelemetry, that is for Cloud Logging(formerly known as Stackdriver Logging)
func initCloudTracer(projectID string, opts ...option.ClientOption) error {
	traceOpts := []cloudtrace.Option{
		cloudtrace.WithTraceClientOptions(opts),
		cloudtrace.WithProjectID(projectID),
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

	props := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
		logzpropagation.HTTPFormat{})
	otel.SetTextMapPropagator(props)

	return nil
}

// httpMiddleware is middleware with net/http instrumentation
func httpMiddleware(label string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {

		return otelhttp.NewHandler(h, label)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logz.Infof(ctx, "log for %s", "info")

		fmt.Fprintf(w, "Hello World")
	})

	// For local project
	mock := cloudmock.NewCloudMock()
	defer mock.Shutdown()
	clientOpts := []option.ClientOption{
		option.WithGRPCConn(mock.ClientConn()),
	}

	projectID := "your project id"

	logz.SetProjectID(projectID)
	if err := initCloudTracer(projectID, clientOpts...); err != nil {
		panic(err)
	}

	h := httpMiddleware("tracer name")(mux)

	log.Fatal(http.ListenAndServe(":8080", h))
}
