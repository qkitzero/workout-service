package set

import (
	"fmt"
	"strings"
)

type Exercise string

func (e Exercise) String() string {
	return string(e)
}

func NewExercise(s string) (Exercise, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Exercise(""), fmt.Errorf("invalid exercise")
	}
	return Exercise(s), nil
}
