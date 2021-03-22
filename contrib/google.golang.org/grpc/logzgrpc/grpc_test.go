package logzgrpc_test

import (
	"context"
	"net"
	"strings"
	"testing"

	"github.com/glassonion1/logz"
	"github.com/glassonion1/logz/contrib/google.golang.org/grpc/logzgrpc"
	"github.com/glassonion1/logz/testhelper"
	"go.opentelemetry.io/otel/trace"
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

func TestUnaryServerInterceptor_integration(t *testing.T) {
	logz.InitTracer()
	logz.SetProjectID("test-project")

	fn := func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		spanCtx := trace.SpanContextFromContext(ctx)
		if spanCtx.TraceID().String() == "00000000000000000000000000000000" {
			t.Error("trace id is zero value")
		}
		if spanCtx.SpanID().String() == "0000000000000000" {
			t.Error("span id is zero value")
		}
		return handler(ctx, req)
	}

	sOpt := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			logzgrpc.UnaryServerInterceptor("test1"),
			fn,
		),
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
	interop.DoEmptyUnaryCall(client)
}

func TestStreamServerInterceptor_integration(t *testing.T) {
	logz.InitTracer()
	logz.SetProjectID("test-project")

	fn := func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		spanCtx := trace.SpanContextFromContext(stream.Context())
		if spanCtx.TraceID().String() == "00000000000000000000000000000000" {
			t.Error("trace id is zero value")
		}
		if spanCtx.SpanID().String() == "0000000000000000" {
			t.Error("span id is zero value")
		}
		return handler(srv, stream)
	}

	sOpt := []grpc.ServerOption{
		grpc.ChainStreamInterceptor(
			logzgrpc.StreamServerInterceptor("test2"),
			fn,
		),
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
	interop.DoServerStreaming(client)
}

func TestInterceptors_AccessLog(t *testing.T) {
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
		got := testhelper.ExtractAccessLogOut(t, func() {
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
		if !strings.Contains(got, `"userAgent":"[grpc-go/1.35.1]"`) {
			t.Error("user agent is not set correctly: error")
		}
	})

	t.Run("StreamServerInterceptor", func(t *testing.T) {
		got := testhelper.ExtractAccessLogOut(t, func() {
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
		if !strings.Contains(got, `"userAgent":"[grpc-go/1.35.1]"`) {
			t.Error("user agent is not set correctly: error")
		}
	})
}
