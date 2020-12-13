package config

import "os"

var (
	ProjectID = ""
)

func init() {
	// In case of App Engine, the value can be obtained.
	// Otherwise, it is an empty string.
	ProjectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
}
