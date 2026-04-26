package muscle

import (
	"context"

	"github.com/qkitzero/workout-service/internal/domain/i18n"
	"github.com/qkitzero/workout-service/internal/domain/muscle"
)

type MuscleUsecase interface {
	ListMuscles(ctx context.Context, lang i18n.Language) ([]muscle.Muscle, error)
}

type muscleUsecase struct {
	muscleRepo muscle.MuscleRepository
}

func NewMuscleUsecase(muscleRepo muscle.MuscleRepository) MuscleUsecase {
	return &muscleUsecase{muscleRepo: muscleRepo}
}

func (u *muscleUsecase) ListMuscles(ctx context.Context, lang i18n.Language) ([]muscle.Muscle, error) {
	return u.muscleRepo.FindAll(ctx, lang)
}
