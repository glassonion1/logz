package config

import (
	"github.com/glassonion1/logz/internal/types"
)

var (
	// ProjectID is gcp project id
	ProjectID = ""
	// WriteAccessLog is function that writes an access log
	WriteAccessLog types.WriteAccessLogFunc
)
