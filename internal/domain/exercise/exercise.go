package exercise

import (
	"github.com/qkitzero/workout-service/internal/domain/muscle"
)

type Exercise interface {
	ID() ExerciseID
	Code() Code
	Category() Category
	Name() Name
	Muscles() []muscle.Muscle
}

type exercise struct {
	id       ExerciseID
	code     Code
	category Category
	name     Name
	muscles  []muscle.Muscle
}

func (e exercise) ID() ExerciseID           { return e.id }
func (e exercise) Code() Code               { return e.code }
func (e exercise) Category() Category       { return e.category }
func (e exercise) Name() Name               { return e.name }
func (e exercise) Muscles() []muscle.Muscle { return e.muscles }

func NewExercise(id ExerciseID, code Code, category Category, name Name, muscles []muscle.Muscle) Exercise {
	return &exercise{
		id:       id,
		code:     code,
		category: category,
		name:     name,
		muscles:  muscles,
	}
}
