package exercise

import (
	"time"

	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/i18n"
	"github.com/qkitzero/workout-service/internal/domain/muscle"
	infraMuscle "github.com/qkitzero/workout-service/internal/infrastructure/muscle"
)

type ExerciseModel struct {
	ID           exercise.ExerciseID
	Code         exercise.Code
	Category     exercise.Category
	CreatedAt    time.Time
	Translations []ExerciseTranslationModel `gorm:"foreignKey:ExerciseID"`
	Muscles      []infraMuscle.MuscleModel  `gorm:"many2many:exercise_muscle;joinForeignKey:exercise_id;joinReferences:muscle_id"`
}

func (ExerciseModel) TableName() string {
	return "exercises"
}

func (m ExerciseModel) ToDomain(lang i18n.Language) exercise.Exercise {
	translations := make([]exercise.Translation, len(m.Translations))
	for i, t := range m.Translations {
		translations[i] = exercise.NewTranslation(t.Lang, t.Name)
	}
	name := exercise.ResolveName(translations, lang, m.Code)

	muscles := make([]muscle.Muscle, len(m.Muscles))
	for i, mg := range m.Muscles {
		muscles[i] = mg.ToDomain(lang)
	}

	return exercise.NewExercise(m.ID, m.Code, m.Category, name, muscles)
}

type ExerciseTranslationModel struct {
	ExerciseID exercise.ExerciseID
	Lang       i18n.Language
	Name       exercise.Name
}

func (ExerciseTranslationModel) TableName() string {
	return "exercise_translations"
}
