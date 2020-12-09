package propagation_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/glassonion1/logz/internal/propagation"
	"github.com/google/go-cmp/cmp"
	"go.opentelemetry.io/otel/oteltest"
	"go.opentelemetry.io/otel/trace"
)

func TestHTTPFormatInject(t *testing.T) {

	mockTracer := oteltest.DefaultTracer()
	ctx, _ := mockTracer.Start(context.Background(), "inject")

	req1 := httptest.NewRequest("GET", "http://example.com", nil)

	// Tests the inject funcation
	hf := propagation.HTTPFormat{}
	hf.Inject(ctx, &req1.Header)

	want := "00000000000000020000000000000000/0000000000000002;o=0"
	got := req1.Header.Get("X-Cloud-Trace-Context")
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("failed to inject test: %v", diff)
	}
}

func TestHTTPFormatExtract(t *testing.T) {

	req1 := httptest.NewRequest("GET", "http://example.com", nil)
	req1.Header.Set("X-Cloud-Trace-Context", "a0d3eee13de6a4bbcf291eb444b94f28/913411593c9338c5;o=1")

	// Tests the extract funcation
	hf := propagation.HTTPFormat{}
	ctx := hf.Extract(context.Background(), req1.Header)

	sc := trace.RemoteSpanContextFromContext(ctx)

	if diff := cmp.Diff(sc.TraceID.String(), "a0d3eee13de6a4bbcf291eb444b94f28"); diff != "" {
		t.Errorf("failed to traceid test: %v", diff)
	}

	if diff := cmp.Diff(sc.SpanID.String(), "913411593c9338c5"); diff != "" {
		t.Errorf("failed to spanid test: %v", diff)
	}

	if sc.TraceFlags != trace.FlagsSampled {
		t.Errorf("failed to trace flag test")
	}
}
