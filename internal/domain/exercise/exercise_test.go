package exercise

import "testing"

func TestNewExercise(t *testing.T) {
	t.Parallel()

	id := NewExerciseID()
	code, err := NewCode("bench_press")
	if err != nil {
		t.Fatalf("failed to new code: %v", err)
	}
	category, err := NewCategory("compound")
	if err != nil {
		t.Fatalf("failed to new category: %v", err)
	}
	jaName, err := NewName("ベンチプレス")
	if err != nil {
		t.Fatalf("failed to new name: %v", err)
	}
	enName, err := NewName("Bench Press")
	if err != nil {
		t.Fatalf("failed to new name: %v", err)
	}

	translations := []Translation{
		NewTranslation(LanguageJa, jaName),
		NewTranslation(Language("en"), enName),
	}

	e := NewExercise(id, code, category, translations)

	if e.ID() != id {
		t.Errorf("ID() = %v, want %v", e.ID(), id)
	}
	if e.Code() != code {
		t.Errorf("Code() = %v, want %v", e.Code(), code)
	}
	if e.Category() != category {
		t.Errorf("Category() = %v, want %v", e.Category(), category)
	}
	if len(e.Translations()) != 2 {
		t.Errorf("len(Translations()) = %d, want %d", len(e.Translations()), 2)
	}
}

func TestExerciseName(t *testing.T) {
	t.Parallel()

	id := NewExerciseID()
	code, _ := NewCode("bench_press")
	category, _ := NewCategory("compound")
	jaName, _ := NewName("ベンチプレス")
	enName, _ := NewName("Bench Press")

	tests := []struct {
		name         string
		translations []Translation
		requestLang  Language
		wantName     string
	}{
		{
			"exact match ja",
			[]Translation{NewTranslation(LanguageJa, jaName), NewTranslation(Language("en"), enName)},
			LanguageJa,
			"ベンチプレス",
		},
		{
			"exact match en",
			[]Translation{NewTranslation(LanguageJa, jaName), NewTranslation(Language("en"), enName)},
			Language("en"),
			"Bench Press",
		},
		{
			"fallback to ja when lang missing",
			[]Translation{NewTranslation(LanguageJa, jaName)},
			Language("en"),
			"ベンチプレス",
		},
		{
			"fallback to first when ja missing",
			[]Translation{NewTranslation(Language("en"), enName)},
			Language("de"),
			"Bench Press",
		},
		{
			"fallback to code when no translations",
			nil,
			LanguageJa,
			"bench_press",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e := NewExercise(id, code, category, tt.translations)
			if got := e.Name(tt.requestLang).String(); got != tt.wantName {
				t.Errorf("Name(%v) = %v, want %v", tt.requestLang, got, tt.wantName)
			}
		})
	}
}
