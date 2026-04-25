package muscle

import (
	"context"

	"github.com/qkitzero/workout-service/internal/domain/i18n"
	"github.com/qkitzero/workout-service/internal/domain/muscle"
)

type MuscleUsecase interface {
	ListMuscles(ctx context.Context, lang string) ([]muscle.Muscle, error)
}

type muscleUsecase struct {
	muscleRepo muscle.MuscleRepository
}

func NewMuscleUsecase(muscleRepo muscle.MuscleRepository) MuscleUsecase {
	return &muscleUsecase{muscleRepo: muscleRepo}
}

func (u *muscleUsecase) ListMuscles(ctx context.Context, lang string) ([]muscle.Muscle, error) {
	language := i18n.LanguageJa
	if lang != "" {
		parsed, err := i18n.NewLanguage(lang)
		if err != nil {
			return nil, err
		}
		language = parsed
	}

	return u.muscleRepo.FindAll(ctx, language)
}
