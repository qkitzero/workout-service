package exercise

import "testing"

func TestNewLanguage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		lang    string
	}{
		{"success ja", true, "ja"},
		{"success en", true, "en"},
		{"success en-US", true, "en-US"},
		{"failure empty", false, ""},
		{"failure uppercase", false, "JA"},
		{"failure too long", false, "japanese"},
		{"failure trailing dash", false, "ja-"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			l, err := NewLanguage(tt.lang)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success && l.String() != tt.lang {
				t.Errorf("String() = %v, want %v", l.String(), tt.lang)
			}
		})
	}
}
