package muscle

import "testing"

func TestNewCode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		code    string
	}{
		{"success simple", true, "chest"},
		{"success with underscore", true, "lower_back"},
		{"success with digit", true, "deltoid_1"},
		{"failure empty", false, ""},
		{"failure uppercase", false, "Chest"},
		{"failure starts with digit", false, "1chest"},
		{"failure hyphen", false, "lower-back"},
		{"failure space", false, "lower back"},
		{"failure non-ascii", false, "胸"},
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
