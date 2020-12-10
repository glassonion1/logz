package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/glassonion1/logz/internal/config"
	"github.com/glassonion1/logz/internal/severity"
	"github.com/glassonion1/logz/internal/spancontext"
)

var NowFunc = time.Now

type LogEntry struct {
	Severity    string       `json:"severity"`
	Message     string       `json:"message"`
	Time        time.Time    `json:"time"`
	Trace       string       `json:"logging.googleapis.com/trace"`
	SpanID      string       `json:"logging.googleapis.com/spanId"`
	HTTPRequest *HttpRequest `json:"httpRequest,omitempty"`
}

type HttpRequest struct {
	RequestMethod string `json:"requestMethod"`
	RequestURL    string `json:"requestUrl"`
	Latency       string `json:"latency"`
	UserAgent     string `json:"userAgent"`
	RemoteIP      string `json:"remoteIp"`
	Status        int32  `json:"status"`
	Protocol      string `json:"protocol"`
	RequestSize   string `json:"requestSize"`
	ResponseSize  string `json:"responseSize"`
}

// Looger is for GCP
type Logger struct {
}

// WriteLog writes a log to stdout
func (l *Logger) WriteLog(ctx context.Context, severity severity.Severity, format string, a ...interface{}) {
	// Gets the traceID and spanID
	sc := spancontext.Extract(ctx)

	trace := fmt.Sprintf("projects/%s/traces/%s", config.ProjectID, sc.TraceID)
	msg := fmt.Sprintf(format, a...)
	ety := &LogEntry{
		Severity: severity.String(),
		Message:  msg,
		Time:     NowFunc(),
		Trace:    trace,
		SpanID:   sc.SpanID,
	}

	if err := json.NewEncoder(os.Stdout).Encode(ety); err != nil {
		panic(err)
	}
}
