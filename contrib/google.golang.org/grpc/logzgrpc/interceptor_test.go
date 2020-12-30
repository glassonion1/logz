package logzgrpc_test

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/glassonion1/logz"
	"github.com/glassonion1/logz/contrib/google.golang.org/grpc/logzgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type stubProtoMessage struct{}

func (s *stubProtoMessage) Reset() {
}

func (s *stubProtoMessage) String() string {
	return "stub"
}

func (s *stubProtoMessage) ProtoMessage() {
}

func TestUnaryServerInterceptor(t *testing.T) {
	logz.InitTracer()
	logz.SetProjectID("test-project")

	usi := logzgrpc.UnaryServerInterceptor("test")
	deniedErr := status.Error(codes.PermissionDenied, "PERMISSION_DENIED_TEXT")
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, deniedErr
	}

	orgStderr := os.Stderr
	defer func() {
		os.Stderr = orgStderr
	}()

	t.Run("Tests unary server interceptor", func(t *testing.T) {
		// Overrides the stderr to the buffer.
		r, w, _ := os.Pipe()
		os.Stderr = w

		_, err := usi(context.Background(), &stubProtoMessage{}, &grpc.UnaryServerInfo{}, handler)

		w.Close()

		var buf bytes.Buffer
		if _, err := buf.ReadFrom(r); err != nil {
			t.Fatalf("failed to read buf: %v", err)
		}
		got := buf.String()

		// Restores the stderr
		os.Stderr = orgStderr

		if err != nil && err.Error() != deniedErr.Error() {
			t.Errorf("unexpected error occured: %s", err)
		}

		if !strings.Contains(got, `"severity":"INFO"`) {
			t.Error("severity is not set correctly: error")
		}

		if !strings.Contains(got, `"logging.googleapis.com/trace":"projects/test-project/traces/`) {
			t.Error("trace is not set correctly: error")
		}

		if !strings.Contains(got, `"httpRequest":{"requestMethod":"gRPC"`) {
			t.Error("http request is not set correctly: error")
		}
	})
}

type stubServerStream struct{}

func (stubServerStream) SetHeader(metadata.MD) error {
	return nil
}

func (stubServerStream) SendHeader(metadata.MD) error {
	return nil
}

func (stubServerStream) SetTrailer(metadata.MD) {}

func (stubServerStream) Context() context.Context {
	return context.Background()
}

func (stubServerStream) SendMsg(m interface{}) error {
	return nil
}

func (stubServerStream) RecvMsg(m interface{}) error {
	return nil
}

func TestStreamServerInterceptor(t *testing.T) {
	logz.InitTracer()
	logz.SetProjectID("test-project")

	ssi := logzgrpc.StreamServerInterceptor("test")
	deniedErr := status.Error(codes.PermissionDenied, "PERMISSION_DENIED_TEXT")
	handler := func(_ interface{}, _ grpc.ServerStream) error {
		return deniedErr
	}

	orgStderr := os.Stderr
	defer func() {
		os.Stderr = orgStderr
	}()

	t.Run("Tests stream server interceptor", func(t *testing.T) {
		// Overrides the stderr to the buffer.
		r, w, _ := os.Pipe()
		os.Stderr = w

		err := ssi(&stubProtoMessage{}, &stubServerStream{}, &grpc.StreamServerInfo{}, handler)

		w.Close()

		var buf bytes.Buffer
		if _, err := buf.ReadFrom(r); err != nil {
			t.Fatalf("failed to read buf: %v", err)
		}
		got := buf.String()

		// Restores the stderr
		os.Stderr = orgStderr

		if err != nil && err.Error() != deniedErr.Error() {
			t.Errorf("unexpected error occured: %s", err)
		}

		if !strings.Contains(got, `"severity":"INFO"`) {
			t.Error("severity is not set correctly: error")
		}

		if !strings.Contains(got, `"logging.googleapis.com/trace":"projects/test-project/traces/`) {
			t.Error("trace is not set correctly: error")
		}

		if !strings.Contains(got, `"httpRequest":{"requestMethod":"gRPC"`) {
			t.Error("http request is not set correctly: error")
		}
	})
}
