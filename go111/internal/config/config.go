package config

import (
	"io"
	"os"
)

var (
	// ProjectID is gcp project id
	ProjectID = ""
	// CallerSkip is the number of stack frames to ascend
	CallerSkip = 0
	// ApplicationLogOut is io.Writer object for application log
	ApplicationLogOut io.Writer
)

func init() {
	// In case of App Engine, the value can be obtained.
	// Otherwise, it is an empty string.
	ProjectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

	ApplicationLogOut = os.Stdout
}
