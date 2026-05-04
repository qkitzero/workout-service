package workout

import (
	"context"
	"time"

	"github.com/qkitzero/workout-service/internal/domain/user"
)

type WorkoutRepository interface {
	Create(ctx context.Context, w Workout) error
	Update(ctx context.Context, w Workout) error
	FindByID(ctx context.Context, id WorkoutID) (Workout, error)
	FindByUserID(ctx context.Context, userID user.UserID, from, to *time.Time) ([]Workout, error)
	Exists(ctx context.Context, id WorkoutID) (bool, error)
}
