package exercise

import (
	"context"

	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/i18n"
)

type ExerciseUsecase interface {
	ListExercises(ctx context.Context, lang i18n.Language) ([]exercise.Exercise, error)
}

type exerciseUsecase struct {
	exerciseRepo exercise.ExerciseRepository
}

func NewExerciseUsecase(exerciseRepo exercise.ExerciseRepository) ExerciseUsecase {
	return &exerciseUsecase{exerciseRepo: exerciseRepo}
}

func (u *exerciseUsecase) ListExercises(ctx context.Context, lang i18n.Language) ([]exercise.Exercise, error) {
	return u.exerciseRepo.FindAll(ctx, lang)
}
