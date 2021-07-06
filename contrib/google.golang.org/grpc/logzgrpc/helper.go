package logzgrpc

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

var statusMap map[codes.Code]int = map[codes.Code]int{
	codes.OK:                 http.StatusOK,
	codes.Canceled:           http.StatusRequestTimeout,
	codes.Unknown:            http.StatusInternalServerError,
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.Unauthenticated:    http.StatusUnauthorized,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
	codes.FailedPrecondition: http.StatusBadRequest,
	codes.Aborted:            http.StatusConflict,
	codes.OutOfRange:         http.StatusBadRequest,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Internal:           http.StatusInternalServerError,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.DataLoss:           http.StatusInternalServerError,
}

func httpStatusFromCode(code codes.Code) int {
	status, ok := statusMap[code]
	if !ok {
		return http.StatusInternalServerError
	}
	return status
}

func binarySize(val interface{}) int {
	if val == nil {
		return 0
	}

	r := reflect.ValueOf(val)
	switch r.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice: // nilable types
		if r.IsNil() {
			return 0
		}
	}

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.EncodeValue(r)
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
