package set

import (
	"context"
	"time"

	"github.com/qkitzero/workout-service/internal/domain/user"
	"github.com/qkitzero/workout-service/internal/domain/workout"
)

type SetRepository interface {
	Create(ctx context.Context, set Set) error
	Update(ctx context.Context, set Set) error
	Delete(ctx context.Context, id SetID) error
	FindByID(ctx context.Context, id SetID) (Set, error)
	FindByUserID(
		ctx context.Context,
		userID user.UserID,
		from, to *time.Time,
		limit int,
		cursorTrainedAt *time.Time,
		cursorSetID *SetID,
	) ([]Set, error)
	FindByWorkoutID(ctx context.Context, workoutID workout.WorkoutID) ([]Set, error)
}
