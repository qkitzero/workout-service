package paging

import (
	"errors"
	"testing"
)

type sampleCursor struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func TestEncodeCursor(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		cursor sampleCursor
	}{
		{"success encode", sampleCursor{A: 42, B: "hello"}},
		{"success encode zero value", sampleCursor{}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			token := EncodeCursor(tt.cursor)
			if token == "" {
				t.Errorf("expected non-empty token")
			}

			decoded, err := DecodeCursor[sampleCursor](token)
			if err != nil {
				t.Errorf("decode failed: %v", err)
			}
			if decoded != tt.cursor {
				t.Errorf("round-trip mismatch: got %+v, want %+v", decoded, tt.cursor)
			}
		})
	}
}

func TestDecodeCursor(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		token   string
	}{
		{"failure invalid base64", false, "!!!not-base64!!!"},
		{"failure invalid json", false, "aGVsbG8"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := DecodeCursor[sampleCursor](tt.token)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if !tt.success && !errors.Is(err, ErrInvalidPageToken) {
				t.Errorf("expected ErrInvalidPageToken, got %v", err)
			}
		})
	}
}
