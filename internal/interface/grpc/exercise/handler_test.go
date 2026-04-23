package exercise

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"

	exercisev1 "github.com/qkitzero/workout-service/gen/go/exercise/v1"
	appexercise "github.com/qkitzero/workout-service/internal/application/exercise"
	"github.com/qkitzero/workout-service/internal/domain/exercise"
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

			code, _ := exercise.NewCode("bench_press")
			category, _ := exercise.NewCategory("compound")
			name, _ := exercise.NewName("ベンチプレス")
			listed := []appexercise.ListedExercise{
				{
					ID:       exercise.NewExerciseID(),
					Code:     code,
					Name:     name,
					Category: category,
				},
			}

			mockUsecase.EXPECT().ListExercises(tt.ctx, tt.lang).Return(listed, tt.listExercisesErr).AnyTimes()

			handler := NewExerciseHandler(mockUsecase)

			req := &exercisev1.ListExercisesRequest{Lang: tt.lang}
			_, err := handler.ListExercises(tt.ctx, req)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
