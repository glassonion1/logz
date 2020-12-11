# logz

![CI for Pull Request](https://github.com/glassonion1/logz/workflows/CI%20for%20Pull%20Request/badge.svg)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue)](https://godoc.org/github.com/glassonion1/logz)
[![GitHub license](https://img.shields.io/github/license/glassonion1/logz)](https://github.com/glassonion1/logz/blob/main/LICENSE)

Package logz provides the structured and hierarchized log with the OpenTelemetry([https://opentelemetry.io](https://opentelemetry.io)).  
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

logz.SetProjectID("your gcp project id")
logz.InitTracer()
// Sets the middleware
h := middleware.NetHTTP("tracer name")(mux)

log.Fatal(http.ListenAndServe(":8080", h))
```

## Examples
See this sample projects for logz detailed usage  
https://github.com/glassonion1/logz/tree/main/example

## How logs are grouped
The logz leverages the grouping feature of GCP Cloud Logging(a.ka. Stackdriver). See following references for more details.
* https://godoc.org/cloud.google.com/go/logging#hdr-Grouping_Logs_by_Request
