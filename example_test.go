package logz_test

import (
	"log"
	"net/http"

	"github.com/glassonion1/logz"
	"github.com/glassonion1/logz/middleware"
)

func Example() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Writes info log
		logz.Infof(ctx, "writes %s log", "info")
	})

	logz.SetProjectID("your gcp project id")
	logz.InitTracer()
	// Sets the middleware
	h := middleware.NetHTTP("tracer name")(mux)

	log.Fatal(http.ListenAndServe(":8080", h))
}

func ExampleConfig() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Writes info log
		logz.Infof(ctx, "writes %s log", "info")
	})

	logz.SetConfig(logz.Config{
		ProjectID:       "your gcp project id",
		WritesAccessLog: false, // Whether or not to write the access log
	})
	logz.InitTracer()
	// Sets the middleware
	h := middleware.NetHTTP("tracer name")(mux)

	log.Fatal(http.ListenAndServe(":8080", h))
}
