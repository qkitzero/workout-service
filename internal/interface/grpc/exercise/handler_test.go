package exercise

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	exercisev1 "github.com/qkitzero/workout-service/gen/go/exercise/v1"
	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/i18n"
	"github.com/qkitzero/workout-service/internal/domain/muscle"
	mocksappexercise "github.com/qkitzero/workout-service/mocks/application/exercise"
)

func TestListExercises(t *testing.T) {
	t.Parallel()

	code, err := exercise.NewCode("bench_press")
	if err != nil {
		t.Fatalf("failed to new code: %v", err)
	}
	category, err := exercise.NewCategory("compound")
	if err != nil {
		t.Fatalf("failed to new category: %v", err)
	}
	name, err := exercise.NewName("ベンチプレス")
	if err != nil {
		t.Fatalf("failed to new name: %v", err)
	}
	chestCode, err := muscle.NewCode("chest")
	if err != nil {
		t.Fatalf("failed to new muscle code: %v", err)
	}
	chestName, err := muscle.NewName("胸")
	if err != nil {
		t.Fatalf("failed to new muscle name: %v", err)
	}
	muscles := []muscle.Muscle{
		muscle.NewMuscle(muscle.NewMuscleID(), chestCode, chestName),
	}
	sample := exercise.NewExercise(exercise.NewExerciseID(), code, category, name, muscles)

	tests := []struct {
		name             string
		ctx              context.Context
		lang             string
		callUsecase      bool
		wantLang         i18n.Language
		listExercisesErr error
		wantCode         codes.Code
	}{
		{"success list exercises", context.Background(), "ja", true, i18n.LanguageJa, nil, codes.OK},
		{"success default lang", context.Background(), "", true, i18n.LanguageJa, nil, codes.OK},
		{"failure invalid lang", context.Background(), "JA", false, "", nil, codes.InvalidArgument},
		{"failure list exercises error", context.Background(), "ja", true, i18n.LanguageJa, fmt.Errorf("list exercises error"), codes.Internal},
		{"failure status preserved", context.Background(), "ja", true, i18n.LanguageJa, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappexercise.NewMockExerciseUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().ListExercises(tt.ctx, tt.wantLang).Return([]exercise.Exercise{sample}, tt.listExercisesErr).Times(1)
			}

			handler := NewExerciseHandler(mockUsecase)

			req := &exercisev1.ListExercisesRequest{Lang: tt.lang}
			_, err := handler.ListExercises(tt.ctx, req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}
