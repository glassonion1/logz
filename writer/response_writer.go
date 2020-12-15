package writer

import (
	"net/http"
	"sync/atomic"
	"time"
)

// ResponseWriter is size counter for http.ResponseWriter
type ResponseWriter struct {
	http.ResponseWriter
	started    time.Time
	size       uint64
	statusCode int
}

// NewResponseWriter creates new ResponseWriter instance
func NewResponseWriter(rw http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: rw,
		started:        time.Now(),
	}
}

// Write returns underlying Write result, while counting data size
func (rw *ResponseWriter) Write(buf []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(buf)
	atomic.AddUint64(&rw.size, uint64(n))
	return n, err
}

// WriteHeader returns underlying WriteHeader
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Elapsed returns since started value
func (rw *ResponseWriter) Elapsed() time.Duration {
	return time.Since(rw.started)
}

// Size returns repsonse size
func (rw *ResponseWriter) Size() int {
	return int(atomic.LoadUint64(&rw.size))
}

// StatusCode returns sent status code
func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}
