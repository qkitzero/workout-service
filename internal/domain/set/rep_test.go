package set

import "testing"

func TestNewRep(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		rep     int32
	}{
		{"success new rep", true, 10},
		{"failure zero rep", false, 0},
		{"failure negative rep", false, -1},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rep, err := NewRep(tt.rep)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && rep.Int32() != tt.rep {
				t.Errorf("Int32() = %v, want %v", rep.Int32(), tt.rep)
			}
		})
	}
}
