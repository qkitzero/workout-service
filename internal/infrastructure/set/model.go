package set

import (
	"time"

	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/set"
	"github.com/qkitzero/workout-service/internal/domain/user"
	"github.com/qkitzero/workout-service/internal/domain/workout"
)

type SetModel struct {
	ID         set.SetID
	UserID     user.UserID
	WorkoutID  workout.WorkoutID
	ExerciseID exercise.ExerciseID
	Rep        set.Rep
	Weight     set.Weight
	TrainedAt  time.Time
	CreatedAt  time.Time
}

func (SetModel) TableName() string {
	return "sets"
}
