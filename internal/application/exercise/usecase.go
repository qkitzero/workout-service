package exercise

import (
	"context"

	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/i18n"
)

type ExerciseUsecase interface {
	ListExercises(ctx context.Context, lang string) ([]exercise.Exercise, error)
}

type exerciseUsecase struct {
	exerciseRepo exercise.ExerciseRepository
}

func NewExerciseUsecase(exerciseRepo exercise.ExerciseRepository) ExerciseUsecase {
	return &exerciseUsecase{exerciseRepo: exerciseRepo}
}

func (u *exerciseUsecase) ListExercises(ctx context.Context, lang string) ([]exercise.Exercise, error) {
	language := i18n.LanguageJa
	if lang != "" {
		parsed, err := i18n.NewLanguage(lang)
		if err != nil {
			return nil, err
		}
		language = parsed
	}

	return u.exerciseRepo.FindAll(ctx, language)
}
