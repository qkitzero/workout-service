package exercise

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/i18n"
	mocksexercise "github.com/qkitzero/workout-service/mocks/domain/exercise"
)

func TestListExercises(t *testing.T) {
	t.Parallel()

	code, err := exercise.NewCode("bench_press")
	if err != nil {
		t.Errorf("failed to new code: %v", err)
	}
	category, err := exercise.NewCategory("compound")
	if err != nil {
		t.Errorf("failed to new category: %v", err)
	}
	name, err := exercise.NewName("ベンチプレス")
	if err != nil {
		t.Errorf("failed to new name: %v", err)
	}
	sample := exercise.NewExercise(
		exercise.NewExerciseID(),
		code,
		category,
		name,
		nil,
	)

	tests := []struct {
		name        string
		success     bool
		ctx         context.Context
		lang        string
		findAllResp []exercise.Exercise
		findAllErr  error
		wantLang    i18n.Language
	}{
		{"success default lang", true, context.Background(), "", []exercise.Exercise{sample}, nil, i18n.LanguageJa},
		{"success ja", true, context.Background(), "ja", []exercise.Exercise{sample}, nil, i18n.LanguageJa},
		{"success en", true, context.Background(), "en", []exercise.Exercise{sample}, nil, i18n.Language("en")},
		{"success empty result", true, context.Background(), "ja", []exercise.Exercise{}, nil, i18n.LanguageJa},
		{"failure invalid lang", false, context.Background(), "JA", nil, nil, ""},
		{"failure find all error", false, context.Background(), "ja", nil, errors.New("find all error"), i18n.LanguageJa},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockExerciseRepository := mocksexercise.NewMockExerciseRepository(ctrl)
			if tt.wantLang != "" {
				mockExerciseRepository.EXPECT().FindAll(gomock.Any(), tt.wantLang).Return(tt.findAllResp, tt.findAllErr).AnyTimes()
			}

			u := NewExerciseUsecase(mockExerciseRepository)

			_, err := u.ListExercises(tt.ctx, tt.lang)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
