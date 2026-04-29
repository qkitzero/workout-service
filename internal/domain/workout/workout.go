package workout

import (
	"time"

	"github.com/qkitzero/workout-service/internal/domain/user"
)

type Workout interface {
	ID() WorkoutID
	UserID() user.UserID
	StartedAt() time.Time
	FinishedAt() *time.Time
	CreatedAt() time.Time
	IsFinished() bool
}

type workout struct {
	id         WorkoutID
	userID     user.UserID
	startedAt  time.Time
	finishedAt *time.Time
	createdAt  time.Time
}

func (w workout) ID() WorkoutID        { return w.id }
func (w workout) UserID() user.UserID  { return w.userID }
func (w workout) StartedAt() time.Time { return w.startedAt }
func (w workout) FinishedAt() *time.Time {
	if w.finishedAt == nil {
		return nil
	}
	t := *w.finishedAt
	return &t
}
func (w workout) CreatedAt() time.Time { return w.createdAt }
func (w workout) IsFinished() bool     { return w.finishedAt != nil }

func NewWorkout(
	id WorkoutID,
	userID user.UserID,
	startedAt time.Time,
	finishedAt *time.Time,
	createdAt time.Time,
) Workout {
	return &workout{
		id:         id,
		userID:     userID,
		startedAt:  startedAt,
		finishedAt: finishedAt,
		createdAt:  createdAt,
	}
}
