package set

import (
	"time"

	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/user"
	"github.com/qkitzero/workout-service/internal/domain/workout"
)

type Set interface {
	ID() SetID
	UserID() user.UserID
	WorkoutID() workout.WorkoutID
	ExerciseID() exercise.ExerciseID
	Rep() Rep
	Weight() Weight
	TrainedAt() time.Time
	CreatedAt() time.Time
}

type set struct {
	id         SetID
	userID     user.UserID
	workoutID  workout.WorkoutID
	exerciseID exercise.ExerciseID
	rep        Rep
	weight     Weight
	trainedAt  time.Time
	createdAt  time.Time
}

func (s set) ID() SetID                       { return s.id }
func (s set) UserID() user.UserID             { return s.userID }
func (s set) WorkoutID() workout.WorkoutID    { return s.workoutID }
func (s set) ExerciseID() exercise.ExerciseID { return s.exerciseID }
func (s set) Rep() Rep                        { return s.rep }
func (s set) Weight() Weight                  { return s.weight }
func (s set) TrainedAt() time.Time            { return s.trainedAt }
func (s set) CreatedAt() time.Time            { return s.createdAt }

func NewSet(
	id SetID,
	userID user.UserID,
	workoutID workout.WorkoutID,
	exerciseID exercise.ExerciseID,
	rep Rep,
	weight Weight,
	trainedAt time.Time,
	createdAt time.Time,
) Set {
	return &set{
		id:         id,
		userID:     userID,
		workoutID:  workoutID,
		exerciseID: exerciseID,
		rep:        rep,
		weight:     weight,
		trainedAt:  trainedAt,
		createdAt:  createdAt,
	}
}
