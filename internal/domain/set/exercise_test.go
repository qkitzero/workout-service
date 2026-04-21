package set

import "testing"

func TestNewExercise(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		success  bool
		exercise string
	}{
		{"success new exercise", true, "bench press"},
		{"failure empty exercise", false, ""},
		{"failure whitespace only exercise", false, "   "},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			exercise, err := NewExercise(tt.exercise)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && exercise.String() != tt.exercise {
				t.Errorf("String() = %v, want %v", exercise.String(), tt.exercise)
			}
		})
	}
}
