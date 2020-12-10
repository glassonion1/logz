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
	"github.com/google/uuid"
)

var NowFunc = time.Now

type LogEntry struct {
	Severity    string       `json:"severity"`
	Message     string       `json:"message"`
	Time        time.Time    `json:"time"`
	Trace       string       `json:"logging.googleapis.com/trace"`
	SpanID      string       `json:"logging.googleapis.com/spanId"`
	HTTPRequest *HttpRequest `json:"httpRequest,omitempty"`

	// InsertID is a unique ID for the log entry. If you provide this field,
	// the logging service considers other log entries in the same log with the
	// same ID as duplicates which can be removed. If omitted, the logging
	// service will generate a unique ID for this log entry. Note that because
	// this client retries RPCs automatically, it is possible (though unlikely)
	// that an Entry without an InsertID will be written more than once.
	InsertID string `json:"logging.googleapis.com/insertId,omitempty"`
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
	id, _ := uuid.NewUUID()
	msg := fmt.Sprintf(format, a...)
	ety := &LogEntry{
		Severity: severity.String(),
		Message:  msg,
		Time:     NowFunc(),
		Trace:    trace,
		SpanID:   sc.SpanID,
		InsertID: id.String(),
	}

	if err := json.NewEncoder(os.Stdout).Encode(ety); err != nil {
		panic(err)
	}
}
