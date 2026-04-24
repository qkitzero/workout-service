package exercise

import "testing"

func TestNewCategory(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		success  bool
		category string
	}{
		{"success compound", true, "compound"},
		{"success isolation", true, "isolation"},
		{"failure empty", false, ""},
		{"failure unknown", false, "cardio"},
		{"failure uppercase", false, "Compound"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c, err := NewCategory(tt.category)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success && c.String() != tt.category {
				t.Errorf("String() = %v, want %v", c.String(), tt.category)
			}
		})
	}
}
