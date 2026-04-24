package exercise

import "testing"

func TestNewName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		input   string
	}{
		{"success japanese", true, "ベンチプレス"},
		{"success english", true, "Bench Press"},
		{"failure empty", false, ""},
		{"failure whitespace only", false, "   "},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n, err := NewName(tt.input)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success && n.String() == "" {
				t.Errorf("expected non-empty name")
			}
		})
	}
}
