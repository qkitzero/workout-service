package exercise

import "testing"

func TestNewCode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		code    string
	}{
		{"success simple", true, "squat"},
		{"success with underscore", true, "bench_press"},
		{"success with digit", true, "row_1"},
		{"failure empty", false, ""},
		{"failure uppercase", false, "BenchPress"},
		{"failure starts with digit", false, "1rm"},
		{"failure hyphen", false, "bench-press"},
		{"failure space", false, "bench press"},
		{"failure non-ascii", false, "ベンチ"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, err := NewCode(tt.code)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success && c.String() != tt.code {
				t.Errorf("String() = %v, want %v", c.String(), tt.code)
			}
		})
	}
}
