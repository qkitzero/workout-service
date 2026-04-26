package exercise

import (
	"testing"

	"github.com/qkitzero/workout-service/internal/domain/i18n"
	"github.com/qkitzero/workout-service/internal/domain/muscle"
)

func TestNewExercise(t *testing.T) {
	t.Parallel()
	id := NewExerciseID()
	code, err := NewCode("bench_press")
	if err != nil {
		t.Errorf("failed to new code: %v", err)
	}
	category, err := NewCategory("compound")
	if err != nil {
		t.Errorf("failed to new category: %v", err)
	}
	name, err := NewName("ベンチプレス")
	if err != nil {
		t.Errorf("failed to new name: %v", err)
	}
	chestID := muscle.NewMuscleID()
	chestCode, err := muscle.NewCode("chest")
	if err != nil {
		t.Errorf("failed to new muscle code: %v", err)
	}
	chestName, err := muscle.NewName("胸")
	if err != nil {
		t.Errorf("failed to new muscle name: %v", err)
	}
	muscles := []muscle.Muscle{muscle.NewMuscle(chestID, chestCode, chestName)}

	tests := []struct {
		name     string
		success  bool
		id       ExerciseID
		code     Code
		category Category
		exName   Name
		muscles  []muscle.Muscle
	}{
		{"success new exercise", true, id, code, category, name, muscles},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := NewExercise(tt.id, tt.code, tt.category, tt.exName, tt.muscles)
			if tt.success && e.ID() != tt.id {
				t.Errorf("ID() = %v, want %v", e.ID(), tt.id)
			}
			if tt.success && e.Code() != tt.code {
				t.Errorf("Code() = %v, want %v", e.Code(), tt.code)
			}
			if tt.success && e.Category() != tt.category {
				t.Errorf("Category() = %v, want %v", e.Category(), tt.category)
			}
			if tt.success && e.Name() != tt.exName {
				t.Errorf("Name() = %v, want %v", e.Name(), tt.exName)
			}
			if tt.success && len(e.Muscles()) != len(tt.muscles) {
				t.Errorf("len(Muscles()) = %d, want %d", len(e.Muscles()), len(tt.muscles))
			}
		})
	}
}

func TestResolveName(t *testing.T) {
	t.Parallel()
	code, err := NewCode("bench_press")
	if err != nil {
		t.Errorf("failed to new code: %v", err)
	}
	jaName, err := NewName("ベンチプレス")
	if err != nil {
		t.Errorf("failed to new name: %v", err)
	}
	enName, err := NewName("Bench Press")
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
			"ベンチプレス",
		},
		{
			"exact match en",
			[]Translation{NewTranslation(i18n.LanguageJa, jaName), NewTranslation(i18n.Language("en"), enName)},
			i18n.Language("en"),
			"Bench Press",
		},
		{
			"fallback to ja when lang missing",
			[]Translation{NewTranslation(i18n.LanguageJa, jaName)},
			i18n.Language("en"),
			"ベンチプレス",
		},
		{
			"fallback to code when ja missing",
			[]Translation{NewTranslation(i18n.Language("en"), enName)},
			i18n.Language("de"),
			"bench_press",
		},
		{
			"fallback to code when no translations",
			nil,
			i18n.LanguageJa,
			"bench_press",
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
