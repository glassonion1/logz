package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/glassonion1/logz/internal/config"
	"github.com/glassonion1/logz/internal/severity"
	"github.com/glassonion1/logz/internal/spancontext"
	"github.com/glassonion1/logz/internal/types"
)

// NowFunc is a function used in testing instead of time.Now
var NowFunc = time.Now

const traceFmt = "projects/%s/traces/%s"

// WriteApplicationLog writes a application log to stdout
func WriteApplicationLog(ctx context.Context, s severity.Severity, format string, a ...interface{}) {
	// Add a severity to the ContextSeverity
	cs := severity.GetContextSeverity(ctx)
	if cs != nil {
		cs.Add(s)
	}

	// Gets the traceID and spanID
	sc := spancontext.Extract(ctx)

	// Gets the source location
	var location types.SourceLocation
	if pc, file, line, ok := runtime.Caller(2 + config.CallerSkip); ok {
		if function := runtime.FuncForPC(pc); function != nil {
			location.Function = function.Name()
		}
		location.Line = fmt.Sprintf("%d", line)
		parts := strings.Split(file, "/")
		location.File = parts[len(parts)-1] // use short file name
	}

	trace := fmt.Sprintf(traceFmt, config.ProjectID, sc.TraceID)
	msg := fmt.Sprintf(format, a...)
	ety := &types.ApplicationLog{
		Severity:       s.String(),
		Message:        msg,
		Time:           NowFunc(),
		SourceLocation: location,
		Trace:          trace,
		SpanID:         sc.SpanID,
		TraceSampled:   sc.TraceSampled,
	}

	w := config.ApplicationLogOut
	if cc := config.GetContextConfig(ctx); cc != nil && cc.ApplicationLogOut != nil {
		w = cc.ApplicationLogOut
	}

	if err := json.NewEncoder(w).Encode(ety); err != nil {
		fmt.Printf("failed to write log: %v", err)
	}
}

// WriteAccessLog writes a access log to stderr
func WriteAccessLog(ctx context.Context, req types.HTTPRequest) {
	// Gets the ContextSeverity
	cs := severity.GetContextSeverity(ctx)
	s := severity.Info
	if req.Status >= 500 {
		s = severity.Error
	} else if req.Status >= 400 {
		s = severity.Warning
	}
	if cs != nil {
		max := cs.Max()
		if max > s {
			s = max
		}
	}

	// Gets the traceID and spanID
	sc := spancontext.Extract(ctx)

	trace := fmt.Sprintf(traceFmt, config.ProjectID, sc.TraceID)
	ety := &types.AccessLog{
		Severity:    s.String(),
		Time:        NowFunc(), // Timestamp marks the end of the request.
		Trace:       trace,
		HTTPRequest: req,
	}

	w := config.AccessLogOut
	if cc := config.GetContextConfig(ctx); cc != nil && cc.AccessLogOut != nil {
		w = cc.AccessLogOut
	}

	if err := json.NewEncoder(w).Encode(ety); err != nil {
		fmt.Printf("failed to write log: %v", err)
	}
}
