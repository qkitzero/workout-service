package muscle

import (
	"fmt"

	"github.com/google/uuid"
)

type MuscleID struct {
	uuid.UUID
}

func NewMuscleID() MuscleID {
	id := uuid.New()
	return MuscleID{id}
}

func NewMuscleIDFromString(s string) (MuscleID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return MuscleID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return MuscleID{id}, nil
}
