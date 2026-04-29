package set

import (
	"context"

	"github.com/qkitzero/workout-service/internal/domain/user"
	"github.com/qkitzero/workout-service/internal/domain/workout"
)

type SetRepository interface {
	Create(ctx context.Context, set Set) error
	FindByUserID(ctx context.Context, userID user.UserID) ([]Set, error)
	FindByWorkoutID(ctx context.Context, workoutID workout.WorkoutID) ([]Set, error)
}
