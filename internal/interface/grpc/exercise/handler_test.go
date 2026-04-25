package exercise

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"

	exercisev1 "github.com/qkitzero/workout-service/gen/go/exercise/v1"
	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/muscle"
	mocksappexercise "github.com/qkitzero/workout-service/mocks/application/exercise"
)

func TestListExercises(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		success          bool
		ctx              context.Context
		lang             string
		listExercisesErr error
	}{
		{"success list exercises", true, context.Background(), "ja", nil},
		{"success default lang", true, context.Background(), "", nil},
		{"failure list exercises error", false, context.Background(), "ja", fmt.Errorf("list exercises error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappexercise.NewMockExerciseUsecase(ctrl)

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

			chestCode, err := muscle.NewCode("chest")
			if err != nil {
				t.Errorf("failed to new muscle code: %v", err)
			}
			chestName, err := muscle.NewName("胸")
			if err != nil {
				t.Errorf("failed to new muscle name: %v", err)
			}
			muscles := []muscle.Muscle{
				muscle.NewMuscle(muscle.NewMuscleID(), chestCode, chestName),
			}

			sample := exercise.NewExercise(
				exercise.NewExerciseID(),
				code,
				category,
				name,
				muscles,
			)

			mockUsecase.EXPECT().ListExercises(tt.ctx, tt.lang).Return([]exercise.Exercise{sample}, tt.listExercisesErr).AnyTimes()

			handler := NewExerciseHandler(mockUsecase)

			req := &exercisev1.ListExercisesRequest{Lang: tt.lang}
			_, err = handler.ListExercises(tt.ctx, req)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
