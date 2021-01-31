package logz_test

import (
	"log"
	"net/http"
	"os"

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

func ExampleSetConfig() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Writes info log
		logz.Infof(ctx, "writes %s log", "info")
	})

	logz.SetConfig(logz.Config{
		ProjectID:      "your gcp project id",
		NeedsAccessLog: false, // Whether or not to write the access log
	})
	logz.InitTracer()
	// Sets the middleware
	h := middleware.NetHTTP("tracer name")(mux)

	log.Fatal(http.ListenAndServe(":8080", h))
}

func ExampleSetConfig_changeLogout() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Writes info log
		logz.Infof(ctx, "writes %s log", "info")
	})

	// Writes log on local file
	file, err := os.OpenFile("local.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	logz.SetConfig(logz.Config{
		ProjectID:         "your gcp project id",
		NeedsAccessLog:    false,
		AccessLogOut:      file,
		ApplicationLogOut: file,
	})
	logz.InitTracer()
	// Sets the middleware
	h := middleware.NetHTTP("tracer name")(mux)

	log.Fatal(http.ListenAndServe(":8080", h))
}
