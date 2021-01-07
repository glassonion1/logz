# logz

![CI for Pull Request](https://github.com/glassonion1/logz/workflows/CI%20for%20Pull%20Request/badge.svg)
[![Godoc](https://img.shields.io/badge/godoc-reference-blue)](https://godoc.org/github.com/glassonion1/logz)
[![Go Report Card](https://goreportcard.com/badge/github.com/glassonion1/logz)](https://goreportcard.com/report/github.com/glassonion1/logz)
[![GitHub license](https://img.shields.io/github/license/glassonion1/logz)](https://github.com/glassonion1/logz/blob/main/LICENSE)

The logz is logger library in Go for grouping application logs related a access log. logz uses OpenTelemetry([https://opentelemetry.io](https://opentelemetry.io)) to generate the trace id.  
This is for Google Cloud Logging (formerly known as Stackdriver Logging).  

![screenshot](https://github.com/glassonion1/logz/blob/main/img/screenshot.png "Cloug Logging")

## Features
* Writes access log each http requests
* Writes application log
* Grouping application logs related a access log.
  * The parent entry will inherit the severity of its children
* Supports to App Engine 2nd and Cloud Run and GKE.

Use [go111](https://github.com/glassonion1/logz/tree/main/go111) package if your project is App Engine 1st generation.

### Contribution Packages
The logz contribution packages that provides middlewares for 3rd-party Go packages.
* [gin](https://github.com/glassonion1/logz/tree/main/contrib/github.com/gin-gonic/gin/logzgin)
* [echo](https://github.com/glassonion1/logz/tree/main/contrib/github.com/labstack/echo/logzecho)
* [gRPC](https://github.com/glassonion1/logz/tree/main/contrib/google.golang.org/grpc/logzgrpc)

For more details: [https://github.com/glassonion1/logz/tree/main/contrib](https://github.com/glassonion1/logz/tree/main/contrib)

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

### Recommended settings
#### GAE
```go
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // Writes info log
        logz.Infof(ctx, "writes %s log", "info")
    })

    logz.SetConfig(logz.Config{
        WritesAccessLog: true, // Writes the access log
    })
    logz.InitTracer()
    // Sets the middleware
    h := middleware.NetHTTP("tracer name")(mux)

    log.Fatal(http.ListenAndServe(":8080", h))
}
```
#### Cloud Run
```go
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // Writes info log
        logz.Infof(ctx, "writes %s log", "info")
    })

    logz.SetConfig(logz.Config{
        ProjectID:       "your gcp project id",
        WritesAccessLog: false, // Writes no access log
    })
    logz.InitTracer()
    // Sets the middleware
    h := middleware.NetHTTP("tracer name")(mux)

    log.Fatal(http.ListenAndServe(":8080", h))
}
```
#### GKE
```go
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // Writes info log
        logz.Infof(ctx, "writes %s log", "info")
    })

    logz.SetConfig(logz.Config{
        ProjectID:       "your gcp project id",
        WritesAccessLog: true, // Writes the access log
    })
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

## Log format
The log format is based on [LogEntry](https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry)'s structured payload

### Application log format

```json
{
  "severity":"INFO",
  "message":"writes info log",
  "time":"2020-12-31T23:59:59.999999999Z",
  "logging.googleapis.com/sourceLocation":{
    "file":"logger_test.go",
    "line":"57",
    "function":"github.com/glassonion1/logz/internal/logger_test.TestLoggerWriteApplicationLog.func3"
  },
  "logging.googleapis.com/trace":"projects/test/traces/00000000000000000000000000000000",
  "logging.googleapis.com/spanId":"0000000000000000",
  "logging.googleapis.com/trace_sampled":false
}
```

### Access log format

```json
{
  "severity":"DEFAULT",
  "time":"2020-12-31T23:59:59.999999999Z",
  "logging.googleapis.com/trace":"projects/test/traces/a0d3eee13de6a4bbcf291eb444b94f28",
  "httpRequest":{
    "requestMethod":"GET",
    "requestUrl":"/test1",
    "requestSize":"0",
    "status":200,
    "responseSize":"333",
    "remoteIp":"192.0.2.1",
    "serverIp":"192.168.100.115",
    "latencyy":{
      "nanos":100,
      "seconds":0
    },
    "protocol":"HTTP/1.1"
  }
}
```
