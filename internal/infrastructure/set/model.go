package set

import (
	"time"

	"github.com/qkitzero/workout-service/internal/domain/set"
	"github.com/qkitzero/workout-service/internal/domain/user"
)

type SetModel struct {
	ID        set.SetID
	UserID    user.UserID
	Exercise  set.Exercise
	Rep       set.Rep
	Weight    set.Weight
	TrainedAt time.Time
	CreatedAt time.Time
}

func (SetModel) TableName() string {
	return "sets"
}
