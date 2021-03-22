package logzgrpc

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/glassonion1/logz"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type metadataSupplier struct {
	metadata *metadata.MD
}

// Get returns the value associated with the passed key.
func (s *metadataSupplier) Get(key string) string {
	values := s.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

// Set stores the key-value pair.
func (s *metadataSupplier) Set(key string, value string) {
	s.metadata.Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (s *metadataSupplier) Keys() []string {
	keys := make([]string, 0, s.metadata.Len())
	md := *s.metadata
	for k := range md {
		keys = append(keys, k)
	}
	return keys
}

// UnaryServerInterceptor returns a grpc.UnaryServerInterceptor suitable
// for use in a grpc.NewServer call.
func UnaryServerInterceptor(label string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		started := time.Now()
		md, _ := metadata.FromIncomingContext(ctx)
		metadataCopy := md.Copy()

		tracer := otel.Tracer(label)
		prop := otel.GetTextMapPropagator()
		ctx = prop.Extract(ctx, &metadataSupplier{
			metadata: &metadataCopy,
		})

		ctx, span := tracer.Start(ctx, info.FullMethod)
		ctx = logz.StartCollectingSeverity(ctx)

		var res interface{}
		var err error
		defer func() {
			ua, ip := extractUAAndIP(md)
			reqSize := binarySize(req)
			resSize := binarySize(res)
			code := httpStatusFromCode(status.Code(err))

			logz.AccessLog(ctx, "gRPC Unary", info.FullMethod,
				ua, ip, "HTTP/2",
				code, reqSize, resSize, time.Since(started))
			span.End()
		}()

		res, err = handler(ctx, req)
		return res, err
	}
}

// serverStream wraps around the embedded grpc.ServerStream, and intercepts the RecvMsg and
// SendMsg method call.
type serverStream struct {
	grpc.ServerStream
	ctx context.Context

	requestSize  uint64
	responseSize uint64
}

func (s *serverStream) Context() context.Context {
	return s.ctx
}

func (s *serverStream) SendMsg(m interface{}) error {
	err := s.ServerStream.SendMsg(m)
	if err == nil {
		atomic.AddUint64(&s.responseSize, uint64(binarySize(m)))
	}
	return err
}

func (s *serverStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if err == nil {
		atomic.AddUint64(&s.requestSize, uint64(binarySize(m)))
	}
	return err
}

// StreamServerInterceptor returns a grpc.StreamServerInterceptor suitable
// for use in a grpc.NewServer call.
func StreamServerInterceptor(label string) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		started := time.Now()
		ctx := stream.Context()
		md, _ := metadata.FromIncomingContext(ctx)
		metadataCopy := md.Copy()

		tracer := otel.Tracer(label)
		prop := otel.GetTextMapPropagator()
		ctx = prop.Extract(ctx, &metadataSupplier{
			metadata: &metadataCopy,
		})

		ctx, span := tracer.Start(ctx, info.FullMethod)
		ctx = logz.StartCollectingSeverity(ctx)

		wrapped := &serverStream{
			ServerStream: stream,
			ctx:          ctx,
		}
		var err error
		defer func() {
			ua, ip := extractUAAndIP(md)
			reqSize := int(wrapped.requestSize)
			resSize := int(wrapped.responseSize)
			code := httpStatusFromCode(status.Code(err))

			logz.AccessLog(ctx, "gRPC Server Streaming", info.FullMethod,
				ua, ip, "HTTP/2",
				code, reqSize, resSize, time.Since(started))
			span.End()
		}()

		err = handler(srv, wrapped)
		return err
	}
}
