package exercise

import (
	"time"

	"github.com/qkitzero/workout-service/internal/domain/exercise"
)

type ExerciseModel struct {
	ID           exercise.ExerciseID
	Code         exercise.Code
	Category     exercise.Category
	CreatedAt    time.Time
	Translations []ExerciseTranslationModel `gorm:"foreignKey:ExerciseID;references:ID"`
}

func (ExerciseModel) TableName() string {
	return "exercises"
}

type ExerciseTranslationModel struct {
	ExerciseID exercise.ExerciseID
	Lang       exercise.Language
	Name       exercise.Name
}

func (ExerciseTranslationModel) TableName() string {
	return "exercise_translations"
}
