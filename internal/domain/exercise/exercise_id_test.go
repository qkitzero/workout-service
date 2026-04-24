package exercise

import "testing"

func TestNewExerciseIDFromString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		s       string
	}{
		{"success valid uuid", true, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"},
		{"failure empty", false, ""},
		{"failure not uuid", false, "not-a-uuid"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id, err := NewExerciseIDFromString(tt.s)
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

func TestNewExerciseID(t *testing.T) {
	t.Parallel()
	id := NewExerciseID()
	if id.String() == "" {
		t.Errorf("expected non-empty id")
	}
}
