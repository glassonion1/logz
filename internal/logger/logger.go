package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/glassonion1/logz/internal/severity"
)

var NowFunc = time.Now

type LogEntry struct {
	Severity    severity.Severity `json:"severity,string"`
	Message     string            `json:"message"`
	Time        time.Time         `json:"time"`
	Trace       string            `json:"logging.googleapis.com/trace"`
	SpanID      string            `json:"logging.googleapis.com/spanId"`
	JSONPayload interface{}       `json:"jsonPayload"`
	HTTPRequest *HttpRequest      `json:"httpRequest,omitempty"`
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

type Logger struct {
}

func (l *Logger) WriteLog(ctx context.Context, severity severity.Severity, format string, a ...interface{}) {
	// TODO: gets the trace and spanId

	msg := fmt.Sprintf(format, a...)
	ety := &LogEntry{
		Severity: severity,
		Message:  msg,
		Time:     NowFunc(),
		Trace:    "",
		SpanID:   "",
	}

	if err := json.NewEncoder(os.Stdout).Encode(ety); err != nil {
		panic(err)
	}
}
