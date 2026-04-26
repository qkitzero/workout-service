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
		lang        i18n.Language
		findAllResp []exercise.Exercise
		findAllErr  error
	}{
		{"success ja", true, context.Background(), i18n.LanguageJa, []exercise.Exercise{sample}, nil},
		{"success en", true, context.Background(), i18n.Language("en"), []exercise.Exercise{sample}, nil},
		{"success empty result", true, context.Background(), i18n.LanguageJa, []exercise.Exercise{}, nil},
		{"failure find all error", false, context.Background(), i18n.LanguageJa, nil, errors.New("find all error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockExerciseRepository := mocksexercise.NewMockExerciseRepository(ctrl)
			mockExerciseRepository.EXPECT().FindAll(gomock.Any(), tt.lang).Return(tt.findAllResp, tt.findAllErr).AnyTimes()

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
