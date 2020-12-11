package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/glassonion1/logz"
	"github.com/glassonion1/logz/middleware"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logz.Infof(ctx, "log for %s", "info")

		fmt.Fprintf(w, "Hello World")
	})

	logz.SetProjectID("your project id")
	logz.InitTracer()

	h := middleware.NetHTTP("tracer name")(mux)

	log.Fatal(http.ListenAndServe(":8080", h))
}
