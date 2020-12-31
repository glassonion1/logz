package logzgrpc_test

import (
	"bytes"
	"context"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/glassonion1/logz"
	"github.com/glassonion1/logz/contrib/google.golang.org/grpc/logzgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/interop"
	pb "google.golang.org/grpc/interop/grpc_testing"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func serve(sOpt []grpc.ServerOption) *grpc.Server {
	lis = bufconn.Listen(bufSize)

	s := grpc.NewServer(sOpt...)
	pb.RegisterTestServiceServer(s, interop.NewTestServer())
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
	return s
}

func extractStderr(t *testing.T, fnc func()) string {
	// Evacuates the stderr
	orgStderr := os.Stderr
	defer func() {
		os.Stderr = orgStderr
	}()

	// Overrides the stderr to the buffer.
	r, w, _ := os.Pipe()
	os.Stderr = w

	fnc()

	w.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to read buf: %v", err)
	}

	return strings.TrimRight(buf.String(), "\n")
}

func TestInterceptors(t *testing.T) {
	logz.InitTracer()
	logz.SetProjectID("test-project")

	sOpt := []grpc.ServerOption{
		grpc.UnaryInterceptor(logzgrpc.UnaryServerInterceptor("test1")),
		grpc.StreamInterceptor(logzgrpc.StreamServerInterceptor("test2")),
	}
	s := serve(sOpt)
	defer s.Stop()

	ctx := context.Background()
	dial := func(context.Context, string) (net.Conn, error) { return lis.Dial() }
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		append([]grpc.DialOption{
			grpc.WithContextDialer(dial),
			grpc.WithInsecure(),
		})...,
	)
	if err != nil {
		t.Errorf("fialed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewTestServiceClient(conn)

	t.Run("UnaryServerInterceptor", func(t *testing.T) {
		got := extractStderr(t, func() {
			interop.DoEmptyUnaryCall(client)
		})
		if !strings.Contains(got, `"severity":"INFO"`) {
			t.Error("severity is not set correctly: error")
		}
		if !strings.Contains(got, `"logging.googleapis.com/trace":"projects/test-project/traces`) {
			t.Error("trace is not set correctly: error")
		}
		if !strings.Contains(got, `"httpRequest":{"requestMethod":"gRPC Unary"`) {
			t.Error("http request is not set correctly: error")
		}
		if !strings.Contains(got, `"userAgent":"[grpc-go/1.34.0]"`) {
			t.Error("user agent is not set correctly: error")
		}
	})

	t.Run("StreamServerInterceptor", func(t *testing.T) {
		got := extractStderr(t, func() {
			interop.DoServerStreaming(client)
		})
		if !strings.Contains(got, `"severity":"INFO"`) {
			t.Error("severity is not set correctly: error")
		}
		if !strings.Contains(got, `"logging.googleapis.com/trace":"projects/test-project/traces`) {
			t.Error("trace is not set correctly: error")
		}
		if !strings.Contains(got, `"httpRequest":{"requestMethod":"gRPC Server Streaming"`) {
			t.Error("http request is not set correctly: error")
		}
		if !strings.Contains(got, `"userAgent":"[grpc-go/1.34.0]"`) {
			t.Error("user agent is not set correctly: error")
		}
	})
}
