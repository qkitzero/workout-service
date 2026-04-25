package muscle

import (
	"testing"

	"github.com/qkitzero/workout-service/internal/domain/i18n"
)

func TestNewMuscle(t *testing.T) {
	t.Parallel()
	id := NewMuscleID()
	code, err := NewCode("chest")
	if err != nil {
		t.Errorf("failed to new code: %v", err)
	}
	name, err := NewName("胸")
	if err != nil {
		t.Errorf("failed to new name: %v", err)
	}

	tests := []struct {
		name    string
		success bool
		id      MuscleID
		code    Code
		mName   Name
	}{
		{"success new muscle", true, id, code, name},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := NewMuscle(tt.id, tt.code, tt.mName)
			if tt.success && m.ID() != tt.id {
				t.Errorf("ID() = %v, want %v", m.ID(), tt.id)
			}
			if tt.success && m.Code() != tt.code {
				t.Errorf("Code() = %v, want %v", m.Code(), tt.code)
			}
			if tt.success && m.Name() != tt.mName {
				t.Errorf("Name() = %v, want %v", m.Name(), tt.mName)
			}
		})
	}
}

func TestResolveName(t *testing.T) {
	t.Parallel()
	code, err := NewCode("chest")
	if err != nil {
		t.Errorf("failed to new code: %v", err)
	}
	jaName, err := NewName("胸")
	if err != nil {
		t.Errorf("failed to new name: %v", err)
	}
	enName, err := NewName("Chest")
	if err != nil {
		t.Errorf("failed to new name: %v", err)
	}

	tests := []struct {
		name         string
		translations []Translation
		requestLang  i18n.Language
		wantName     string
	}{
		{
			"exact match ja",
			[]Translation{NewTranslation(i18n.LanguageJa, jaName), NewTranslation(i18n.Language("en"), enName)},
			i18n.LanguageJa,
			"胸",
		},
		{
			"exact match en",
			[]Translation{NewTranslation(i18n.LanguageJa, jaName), NewTranslation(i18n.Language("en"), enName)},
			i18n.Language("en"),
			"Chest",
		},
		{
			"fallback to ja when lang missing",
			[]Translation{NewTranslation(i18n.LanguageJa, jaName)},
			i18n.Language("en"),
			"胸",
		},
		{
			"fallback to code when ja missing",
			[]Translation{NewTranslation(i18n.Language("en"), enName)},
			i18n.Language("de"),
			"chest",
		},
		{
			"fallback to code when no translations",
			nil,
			i18n.LanguageJa,
			"chest",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := ResolveName(tt.translations, tt.requestLang, code).String(); got != tt.wantName {
				t.Errorf("ResolveName(%v) = %v, want %v", tt.requestLang, got, tt.wantName)
			}
		})
	}
}
