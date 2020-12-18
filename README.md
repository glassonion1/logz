# logz

![CI for Pull Request](https://github.com/glassonion1/logz/workflows/CI%20for%20Pull%20Request/badge.svg)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue)](https://godoc.org/github.com/glassonion1/logz)
[![Go Report Card](https://goreportcard.com/badge/github.com/glassonion1/logz)](https://goreportcard.com/report/github.com/glassonion1/logz)
[![GitHub license](https://img.shields.io/github/license/glassonion1/logz)](https://github.com/glassonion1/logz/blob/main/LICENSE)

The logz is logger library in Go for grouping application logs related a access log. logz uses OpenTelemetry([https://opentelemetry.io](https://opentelemetry.io)) to generate the trace id.  
This is for Google Cloud Logging (formerly known as Stackdriver Logging).  

## Features
- Writes access log each http requests
- Writes application log
- Grouping application logs related a access log.
- Supports to App Engine 2nd and Cloud Run and GKE.

Use [go111](https://github.com/glassonion1/logz/tree/main/go111) package if your project is App Engine 1st generation.
  
## Install
```
$ go get github.com/glassonion1/logz
```

## Usage

```go
package main

import (
    "log"
    "net/http"

    "github.com/glassonion1/logz"
    "github.com/glassonion1/logz/middleware"
)

func main() {
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
```

## Examples
See this sample projects for logz detailed usage  
https://github.com/glassonion1/logz/tree/main/example

## How logs are grouped
The logz leverages the grouping feature of GCP Cloud Logging. See following references for more details.
* https://godoc.org/cloud.google.com/go/logging#hdr-Grouping_Logs_by_Request
