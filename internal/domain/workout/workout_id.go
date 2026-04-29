package workout

import (
	"fmt"

	"github.com/google/uuid"
)

type WorkoutID struct {
	uuid.UUID
}

func NewWorkoutID() WorkoutID {
	id := uuid.New()
	return WorkoutID{id}
}

func NewWorkoutIDFromString(s string) (WorkoutID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return WorkoutID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return WorkoutID{id}, nil
}
