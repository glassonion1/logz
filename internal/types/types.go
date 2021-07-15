package types

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

// WriteAccessLogFunc is function type that writes an access log
type WriteAccessLogFunc func(ctx context.Context, req HTTPRequest)

// WriteEmptyAccessLog writes empty log
var WriteEmptyAccessLog = func(context.Context, HTTPRequest) {}

// ApplicationLog is a structured log written by the developer by any timing
type ApplicationLog struct {
	Severity       string         `json:"severity"`
	Message        interface{}    `json:"message"`
	Time           time.Time      `json:"time"`
	SourceLocation SourceLocation `json:"logging.googleapis.com/sourceLocation"`
	Trace          string         `json:"logging.googleapis.com/trace"`
	SpanID         string         `json:"logging.googleapis.com/spanId"`
	TraceSampled   bool           `json:"logging.googleapis.com/trace_sampled"`
}

// AccessLog is a log written by the service each time it is accessed by the client
type AccessLog struct {
	Severity    string      `json:"severity"`
	Time        time.Time   `json:"time"`
	Trace       string      `json:"logging.googleapis.com/trace"`
	HTTPRequest HTTPRequest `json:"httpRequest"`
}

// SourceLocation is a location of source
type SourceLocation struct {
	File     string `json:"file"`
	Line     string `json:"line"`
	Function string `json:"function"`
}

// HTTPRequest is a http request struct for the access log
type HTTPRequest struct {
	RequestMethod string   `json:"requestMethod"`
	RequestURL    string   `json:"requestUrl"`
	RequestSize   string   `json:"requestSize,omitempty"`
	Status        int      `json:"status"`
	ResponseSize  string   `json:"responseSize,omitempty"`
	UserAgent     string   `json:"userAgent,omitempty"`
	RemoteIP      string   `json:"remoteIp,omitempty"`
	ServerIP      string   `json:"serverIp,omitempty"`
	Referer       string   `json:"referer,omitempty"`
	Latency       Duration `json:"latency"`
	Protocol      string   `json:"protocol"`
}

// Duration is duration format for protobuf
// see: https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#duration
type Duration struct {
	Nanos   int32 `json:"nanos"`
	Seconds int64 `json:"seconds"`
}

// MakeDuration makes duration struct from time.Duration
func MakeDuration(d time.Duration) Duration {
	nanos := d.Nanoseconds()
	secs := nanos / 1e9
	nanos -= secs * 1e9
	return Duration{
		Nanos:   int32(nanos),
		Seconds: secs,
	}
}

// MakeHTTPRequest makes HTTPRequest struct from http.Request
func MakeHTTPRequest(r http.Request, status, responseSize int, elapsed time.Duration) HTTPRequest {
	return HTTPRequest{
		RequestMethod: r.Method,
		RequestURL:    r.URL.RequestURI(),
		RequestSize:   fmt.Sprintf("%d", r.ContentLength),
		Status:        status,
		ResponseSize:  fmt.Sprintf("%d", responseSize),
		UserAgent:     r.UserAgent(),
		RemoteIP:      strings.Split(r.RemoteAddr, ":")[0],
		ServerIP:      GetServerIP(),
		Referer:       r.Referer(),
		Latency:       MakeDuration(elapsed),
		Protocol:      r.Proto,
	}
}

// GetServerIP is a function for the unit test
var GetServerIP = getServerIP

func getServerIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return ""
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}
			}
		}
	}
	return ""
}
