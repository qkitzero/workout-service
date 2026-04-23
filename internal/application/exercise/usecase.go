package exercise

import (
	"context"

	"github.com/qkitzero/workout-service/internal/domain/exercise"
)

type ListedExercise struct {
	ID       exercise.ExerciseID
	Code     exercise.Code
	Name     exercise.Name
	Category exercise.Category
}

type ExerciseUsecase interface {
	ListExercises(ctx context.Context, lang string) ([]ListedExercise, error)
}

type exerciseUsecase struct {
	exerciseRepo exercise.ExerciseRepository
}

func NewExerciseUsecase(exerciseRepo exercise.ExerciseRepository) ExerciseUsecase {
	return &exerciseUsecase{exerciseRepo: exerciseRepo}
}

func (u *exerciseUsecase) ListExercises(ctx context.Context, lang string) ([]ListedExercise, error) {
	language := exercise.LanguageJa
	if lang != "" {
		parsed, err := exercise.NewLanguage(lang)
		if err != nil {
			return nil, err
		}
		language = parsed
	}

	exercises, err := u.exerciseRepo.FindAll()
	if err != nil {
		return nil, err
	}

	results := make([]ListedExercise, len(exercises))
	for i, e := range exercises {
		results[i] = ListedExercise{
			ID:       e.ID(),
			Code:     e.Code(),
			Name:     e.Name(language),
			Category: e.Category(),
		}
	}
	return results, nil
}
