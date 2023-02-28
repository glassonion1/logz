package config

import (
	"context"
	"io"
	"os"

	"github.com/glassonion1/logz/internal/types"
)

var (
	// ProjectID is gcp project id
	ProjectID = ""
	// CallerDepth is the number of stack frames to ascend
	CallerSkip = 0
	// WriteAccessLog is function that writes an access log
	WriteAccessLog types.WriteAccessLogFunc
	// ApplicationLogOut is io.Writer object for application log
	ApplicationLogOut io.Writer
	// AccessLogOut is io.Writer object for access log
	AccessLogOut io.Writer
)

func init() {
	// In case of App Engine, the value can be obtained.
	// Otherwise, it is an empty string.
	ProjectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

	ApplicationLogOut = os.Stdout
	AccessLogOut = os.Stderr
}

type ContextConfig struct {
	// ApplicationLogOut is io.Writer object for application log
	ApplicationLogOut io.Writer
	// AccessLogOut is io.Writer object for access log
	AccessLogOut io.Writer
}

type contextKey struct{}

var contextConfigKey = &contextKey{}

// GetContextConfig sets the ContextConfig instance to context
func SetContextConfig(ctx context.Context, cs *ContextConfig) context.Context {
	return context.WithValue(ctx, contextConfigKey, cs)
}

// GetContextConfig gets the ContextSeverity instance from context
func GetContextConfig(ctx context.Context) *ContextConfig {
	v, _ := ctx.Value(contextConfigKey).(*ContextConfig)
	return v
}
