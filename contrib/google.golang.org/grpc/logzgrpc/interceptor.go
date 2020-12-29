package logzgrpc

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"strings"
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

func (s *metadataSupplier) Get(key string) string {
	values := s.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (s *metadataSupplier) Set(key string, value string) {
	s.metadata.Set(key, value)
}

func binarySize(val interface{}) int {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(val)
	if err != nil {
		return 0
	}
	return binary.Size(buff.Bytes())
}

func extractUAAndIP(md metadata.MD) (string, string) {
	var ua string
	if v, ok := md["x-forwarded-user-agent"]; ok {
		ua = fmt.Sprintf("%v", v)
	} else {
		ua = fmt.Sprintf("%v", md["user-agent"])
	}

	var ip string
	if v, ok := md["x-forwarded-for"]; ok && len(v) > 0 {
		ips := strings.Split(v[0], ",")
		ip = ips[0]
	}
	return ua, ip
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
			code := int(status.Code(err))

			logz.AccessLog(ctx, "gRPC", info.FullMethod,
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
	requestSize  int
	responseSize int
}

func (s *serverStream) SendMsg(m interface{}) error {
	err := s.ServerStream.SendMsg(m)
	if err == nil {
		s.responseSize += binarySize(m)
	}
	return err
}

func (s *serverStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if err == nil {
		s.requestSize += binarySize(m)
	}
	return err
}

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
		}
		var err error
		defer func() {
			ua, ip := extractUAAndIP(md)
			reqSize := wrapped.requestSize
			resSize := wrapped.responseSize
			code := int(status.Code(err))

			logz.AccessLog(ctx, "gRPC", info.FullMethod,
				ua, ip, "HTTP/2",
				code, reqSize, resSize, time.Since(started))
			span.End()
		}()

		err = handler(srv, wrapped)
		return err
	}
}
