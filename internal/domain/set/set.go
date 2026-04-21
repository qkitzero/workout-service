package set

import (
	"time"

	"github.com/qkitzero/workout-service/internal/domain/user"
)

type Set interface {
	ID() SetID
	UserID() user.UserID
	Exercise() Exercise
	Rep() Rep
	Weight() Weight
	TrainedAt() time.Time
	CreatedAt() time.Time
}

type set struct {
	id        SetID
	userID    user.UserID
	exercise  Exercise
	rep       Rep
	weight    Weight
	trainedAt time.Time
	createdAt time.Time
}

func (s set) ID() SetID            { return s.id }
func (s set) UserID() user.UserID  { return s.userID }
func (s set) Exercise() Exercise   { return s.exercise }
func (s set) Rep() Rep             { return s.rep }
func (s set) Weight() Weight       { return s.weight }
func (s set) TrainedAt() time.Time { return s.trainedAt }
func (s set) CreatedAt() time.Time { return s.createdAt }

func NewSet(
	id SetID,
	userID user.UserID,
	exercise Exercise,
	rep Rep,
	weight Weight,
	trainedAt time.Time,
	createdAt time.Time,
) Set {
	return &set{
		id:        id,
		userID:    userID,
		exercise:  exercise,
		rep:       rep,
		weight:    weight,
		trainedAt: trainedAt,
		createdAt: createdAt,
	}
}
