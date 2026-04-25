package muscle

import (
	"context"

	"github.com/qkitzero/workout-service/internal/domain/i18n"
)

type MuscleRepository interface {
	FindAll(ctx context.Context, lang i18n.Language) ([]Muscle, error)
}
