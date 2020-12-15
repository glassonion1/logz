# logz for Go 1.11
This is for migrating `google.golang.org/appengine/log` of Ape Enigne 1st generation.  
Guaranteed to work on Google App Engine 1st generation of Go 1.11.

![screenshot](https://github.com/glassonion1/logz/blob/main/go111/img/screenshot.png "Cloug Logging")

## Install
```
$ go get github.com/glassonion1/logz/go111
```

## Usage

```go
package main

import (
	"log"
	"net/http"

	logz "github.com/glassonion1/logz/go111"
	"github.com/glassonion1/logz/go111/middleware"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        // Writes info log
        logz.Infof(ctx, "logging...")
    })

    logz.SetProjectID("your gcp project id")
    // Sets the middleware
    h := middleware.NetHTTP("tracer name")(mux)

    log.Fatal(http.ListenAndServe(":8080", h))
}
```

## Migrates to logz from google.golang.org/appengine/log

From App Engine Logger
```go
import "google.golang.org/appengine/log"

func Handler(w http.ResponseWriter, r *http.Request) {
    ctx := appengine.NewContext(r)
    log.Infof(ctx, "write log %v", "info")
}
```

To logz
```go
import(
-   "google.golang.org/appengine/log"
+   log "github.com/glassonion1/logz/go111"
) 

func Handler(w http.ResponseWriter, r *http.Request) {
-   ctx := appengine.NewContext(r)
+   ctx := r.Context()
    log.Infof(ctx, "write log %v", "info")
}
```
Removes appengine.Context to context.Context
