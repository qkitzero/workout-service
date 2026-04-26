package exercise

import (
	"context"

	"github.com/qkitzero/workout-service/internal/domain/i18n"
)

type ExerciseRepository interface {
	FindAll(ctx context.Context, lang i18n.Language) ([]Exercise, error)
	FindByID(ctx context.Context, id ExerciseID, lang i18n.Language) (Exercise, error)
	Exists(ctx context.Context, id ExerciseID) (bool, error)
}
