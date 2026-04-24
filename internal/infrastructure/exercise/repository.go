package exercise

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/qkitzero/workout-service/internal/domain/exercise"
)

type exerciseRepository struct {
	db *gorm.DB
}

func NewExerciseRepository(db *gorm.DB) exercise.ExerciseRepository {
	return &exerciseRepository{db: db}
}

func (r *exerciseRepository) FindAll(ctx context.Context) ([]exercise.Exercise, error) {
	var models []ExerciseModel
	if err := r.db.WithContext(ctx).Preload("Translations").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]exercise.Exercise, len(models))
	for i, m := range models {
		result[i] = toDomain(m)
	}
	return result, nil
}

func (r *exerciseRepository) FindByID(ctx context.Context, id exercise.ExerciseID) (exercise.Exercise, error) {
	var model ExerciseModel
	if err := r.db.WithContext(ctx).Preload("Translations").Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exercise.ErrExerciseNotFound
		}
		return nil, err
	}
	return toDomain(model), nil
}

func (r *exerciseRepository) Exists(ctx context.Context, id exercise.ExerciseID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&ExerciseModel{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func toDomain(m ExerciseModel) exercise.Exercise {
	translations := make([]exercise.Translation, len(m.Translations))
	for i, t := range m.Translations {
		translations[i] = exercise.NewTranslation(t.Lang, t.Name)
	}
	return exercise.NewExercise(m.ID, m.Code, m.Category, translations)
}
