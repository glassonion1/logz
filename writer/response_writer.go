package writer

import (
	"net/http"
	"sync/atomic"
)

// ResponseWriter is size counter for http.ResponseWriter
type ResponseWriter struct {
	http.ResponseWriter
	size       uint64
	statusCode int
}

// NewResponseWriter creates new ResponseWriter instance
func NewResponseWriter(rw http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: rw,
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

// Size returns response size
func (rw *ResponseWriter) Size() int {
	return int(atomic.LoadUint64(&rw.size))
}

// StatusCode returns sent status code
func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}
