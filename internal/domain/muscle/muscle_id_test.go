package muscle

import "testing"

func TestNewMuscleIDFromString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		s       string
	}{
		{"success valid uuid", true, "4b5a784a-3333-4721-a071-2e3fbd570c7f"},
		{"failure empty", false, ""},
		{"failure not uuid", false, "not-a-uuid"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id, err := NewMuscleIDFromString(tt.s)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success && id.String() != tt.s {
				t.Errorf("String() = %v, want %v", id.String(), tt.s)
			}
		})
	}
}

func TestNewMuscleID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
	}{
		{"success new muscle id", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id := NewMuscleID()
			if tt.success && id.String() == "" {
				t.Errorf("expected non-empty id")
			}
		})
	}
}
