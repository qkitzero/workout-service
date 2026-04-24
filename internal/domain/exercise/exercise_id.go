package exercise

import (
	"fmt"

	"github.com/google/uuid"
)

type ExerciseID struct {
	uuid.UUID
}

func NewExerciseID() ExerciseID {
	id := uuid.New()
	return ExerciseID{id}
}

func NewExerciseIDFromString(s string) (ExerciseID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return ExerciseID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return ExerciseID{id}, nil
}
