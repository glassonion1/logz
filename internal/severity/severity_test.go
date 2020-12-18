package severity_test

import (
	"testing"

	"github.com/glassonion1/logz/internal/severity"
	"github.com/google/go-cmp/cmp"
)

func TestContextSeverity(t *testing.T) {
	tests := []struct {
		name string
		in   []severity.Severity
		want severity.Severity
	}{
		{
			name: "test max severity",
			in:   []severity.Severity{},
			want: severity.Default,
		},
		{
			name: "test max severity",
			in: []severity.Severity{
				severity.Default,
				severity.Error,
			},
			want: severity.Error,
		},
		{
			name: "test max severity",
			in: []severity.Severity{
				severity.Default,
				severity.Error,
				severity.Info,
			},
			want: severity.Error,
		},
		{
			name: "test max severity",
			in: []severity.Severity{
				severity.Default,
				severity.Error,
				severity.Info,
				severity.Critical,
			},
			want: severity.Critical,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cs := severity.ContextSeverity{}
			for _, s := range tt.in {
				cs.Add(s)
			}

			got := cs.Max()

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("failed test %s: %v", tt.name, diff)
			}
		})
	}
}
