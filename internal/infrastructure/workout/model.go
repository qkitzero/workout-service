package workout

import (
	"time"

	"github.com/qkitzero/workout-service/internal/domain/user"
	"github.com/qkitzero/workout-service/internal/domain/workout"
)

type WorkoutModel struct {
	ID         workout.WorkoutID
	UserID     user.UserID
	StartedAt  time.Time
	FinishedAt *time.Time
	CreatedAt  time.Time
}

func (WorkoutModel) TableName() string {
	return "workouts"
}
