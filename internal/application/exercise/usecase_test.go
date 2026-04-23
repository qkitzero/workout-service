package exercise

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/qkitzero/workout-service/internal/domain/exercise"
	mocksexercise "github.com/qkitzero/workout-service/mocks/domain/exercise"
)

func TestListExercises(t *testing.T) {
	t.Parallel()

	code, _ := exercise.NewCode("bench_press")
	category, _ := exercise.NewCategory("compound")
	jaName, _ := exercise.NewName("ベンチプレス")
	enName, _ := exercise.NewName("Bench Press")
	sample := exercise.NewExercise(
		exercise.NewExerciseID(),
		code,
		category,
		[]exercise.Translation{
			exercise.NewTranslation(exercise.LanguageJa, jaName),
			exercise.NewTranslation(exercise.Language("en"), enName),
		},
	)

	tests := []struct {
		name        string
		success     bool
		ctx         context.Context
		lang        string
		findAllResp []exercise.Exercise
		findAllErr  error
		wantName    string
	}{
		{"success default lang", true, context.Background(), "", []exercise.Exercise{sample}, nil, "ベンチプレス"},
		{"success ja", true, context.Background(), "ja", []exercise.Exercise{sample}, nil, "ベンチプレス"},
		{"success en", true, context.Background(), "en", []exercise.Exercise{sample}, nil, "Bench Press"},
		{"success empty result", true, context.Background(), "ja", []exercise.Exercise{}, nil, ""},
		{"failure invalid lang", false, context.Background(), "JA", nil, nil, ""},
		{"failure find all error", false, context.Background(), "ja", nil, errors.New("find all error"), ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockExerciseRepository := mocksexercise.NewMockExerciseRepository(ctrl)
			mockExerciseRepository.EXPECT().FindAll().Return(tt.findAllResp, tt.findAllErr).AnyTimes()

			u := NewExerciseUsecase(mockExerciseRepository)

			got, err := u.ListExercises(tt.ctx, tt.lang)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success && tt.wantName != "" {
				if len(got) != 1 {
					t.Fatalf("expected 1 result, got %d", len(got))
				}
				if got[0].Name.String() != tt.wantName {
					t.Errorf("Name = %v, want %v", got[0].Name.String(), tt.wantName)
				}
			}
		})
	}
}
