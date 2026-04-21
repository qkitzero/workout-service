package set

import (
	"fmt"

	"github.com/google/uuid"
)

type SetID struct {
	uuid.UUID
}

func NewSetID() SetID {
	id := uuid.New()
	return SetID{id}
}

func NewSetIDFromString(s string) (SetID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return SetID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return SetID{id}, nil
}
