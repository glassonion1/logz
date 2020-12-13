package types

import (
	"time"
)

// ApplicationLog is a log written by the developer by any timing
type ApplicationLog struct {
	Severity       string         `json:"severity"`
	Message        string         `json:"message"`
	Time           time.Time      `json:"time"`
	SourceLocation SourceLocation `json:"logging.googleapis.com/sourceLocation"`
	Trace          string         `json:"logging.googleapis.com/trace"`
	SpanID         string         `json:"logging.googleapis.com/spanId"`
	TraceSampled   bool           `json:"logging.googleapis.com/trace_sampled"`
}

type SourceLocation struct {
	File     string `json:"file"`
	Line     string `json:"line"`
	Function string `json:"function"`
}
