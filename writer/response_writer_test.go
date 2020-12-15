package writer_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glassonion1/logz/writer"
)

func TestResponseWriter(t *testing.T) {
	data := "Hello, World!"
	dataLen := len([]byte(data))

	mux := http.NewServeMux()
	mux.Handle("/test1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, data)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test1", nil)
	w := httptest.NewRecorder()
	rw := writer.NewResponseWriter(w)

	mux.ServeHTTP(rw, req)

	if int64(rw.Elapsed()) == 0 {
		t.Fatalf("time measurement is not working")
	}
	if rw.Size() != dataLen {
		t.Fatalf("size mismatch len of test data: %d != %d", rw.Size(), dataLen)
	}
	if rw.StatusCode() != http.StatusOK {
		t.Fatalf("status code mismatch: %d != %d", rw.StatusCode(), http.StatusOK)
	}
}
