# logz

[![Godoc](https://img.shields.io/badge/godoc-reference-blue)](https://godoc.org/github.com/glassonion1/logz)
[![GitHub license](https://img.shields.io/github/license/glassonion1/logz)](https://github.com/glassonion1/logz/blob/main/LICENSE)

Package logz provides the structured log with the OpenTelemetry([OpenTelemetry](https://opentelemetry.io)).  
This is for Google Cloud Logging (formerly known as Stackdriver).

## Install
```
$ go get github.com/glassonion1/logz
```

## Usage

```go
mux := http.NewServeMux()
mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    logz.Infof(ctx, "logging...")
})

h := logz.HTTPMiddleware("tracer name")(mux)

log.Fatal(http.ListenAndServe(":8080", h))
```
