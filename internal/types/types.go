package types

import (
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

type SourceLocation struct {
	File     string `json:"file"`
	Line     string `json:"line"`
	Function string `json:"function"`
}

type HTTPRequest struct {
	RequestMethod string `json:"requestMethod"`
	RequestURL    string `json:"requestUrl"`
	RequestSize   string `json:"requestSize"`
	Status        int    `json:"status"`
	ResponseSize  string `json:"responseSize"`
	UserAgent     string `json:"userAgent"`
	RemoteIP      string `json:"remoteIp"`
	//ServerIP      string `json:"serverIp"`
	//Referer       string `json:"referer"`
	Latency  string `json:"latency"`
	Protocol string `json:"protocol"`
}

func MakeHTTPRequest(r http.Request, status, responseSize int, elapsed time.Duration) HTTPRequest {
	return HTTPRequest{
		RequestMethod: r.Method,
		RequestURL:    r.URL.RequestURI(),
		RequestSize:   fmt.Sprintf("%d", r.ContentLength),
		Status:        status,
		ResponseSize:  fmt.Sprintf("%d", responseSize),
		UserAgent:     r.UserAgent(),
		RemoteIP:      strings.Split(r.RemoteAddr, ":")[0],
		//		ServerIP:      getServerIp(),
		//		Referer:       r.Referer(),
		Latency:  fmt.Sprintf("%fs", elapsed.Seconds()),
		Protocol: r.Proto,
	}
}

func getServerIp() string {
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
