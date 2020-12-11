# logz

![CI for Pull Request](https://github.com/glassonion1/logz/workflows/CI%20for%20Pull%20Request/badge.svg)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue)](https://godoc.org/github.com/glassonion1/logz)
[![GitHub license](https://img.shields.io/github/license/glassonion1/logz)](https://github.com/glassonion1/logz/blob/main/LICENSE)

Package logz provides the structured log with the OpenTelemetry([https://opentelemetry.io](https://opentelemetry.io)).  
This is for Google Cloud Logging (formerly known as Stackdriver).  
The logz supports to App Engine and Cloud Run and GKE.

## Install
```
$ go get github.com/glassonion1/logz
```

## Usage

```go
mux := http.NewServeMux()
mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Writes info log
    logz.Infof(ctx, "logging...")
})

logz.SetProjectID("your project id")
logz.InitTracer()

h := middleware.NetHTTP("tracer name")(mux)

log.Fatal(http.ListenAndServe(":8080", h))
```
