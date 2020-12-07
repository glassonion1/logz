/*
Package logz provides the structured log with the OpenTelemetry.

Example:
	ctx := r.Context() // r is *http.Request
	logz.Infof(ctx, "info log. requestURL: %s", r.URL.String())
*/
package logz // import "github.com/glassonion1/logz"
