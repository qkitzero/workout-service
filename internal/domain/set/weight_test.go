package set

import "testing"

func TestNewWeight(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		weight  float64
	}{
		{"success new weight", true, 60.5},
		{"success zero weight", true, 0},
		{"failure negative weight", false, -1},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			weight, err := NewWeight(tt.weight)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && weight.Float64() != tt.weight {
				t.Errorf("Float64() = %v, want %v", weight.Float64(), tt.weight)
			}
		})
	}
}
