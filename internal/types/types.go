package types

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
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
	Latency       Duration `json:"latencyy"`
	Protocol      string   `json:"protocol"`
}

var _ json.Marshaler = Duration(0)

// Duration provides time.Duration by protobuf format.
type Duration time.Duration

// MarshalJSON convert raw value to JSON value.
func (d Duration) MarshalJSON() ([]byte, error) {
	nanos := time.Duration(d).Nanoseconds()
	secs := nanos / 1e9
	nanos -= secs * 1e9

	v := make(map[string]interface{})
	v["seconds"] = int64(secs)
	v["nanos"] = int32(nanos)

	return json.Marshal(v)
}

// MakeHTTPRequest makes HTTPRequest struct.
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
		Latency:       Duration(elapsed),
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
