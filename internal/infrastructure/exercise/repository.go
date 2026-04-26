package exercise

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/i18n"
)

type exerciseRepository struct {
	db *gorm.DB
}

func NewExerciseRepository(db *gorm.DB) exercise.ExerciseRepository {
	return &exerciseRepository{db: db}
}

func (r *exerciseRepository) FindAll(ctx context.Context, lang i18n.Language) ([]exercise.Exercise, error) {
	var models []ExerciseModel
	if err := r.db.WithContext(ctx).Preload("Translations").Preload("Muscles.Translations").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]exercise.Exercise, len(models))
	for i, m := range models {
		result[i] = m.ToDomain(lang)
	}
	return result, nil
}

func (r *exerciseRepository) FindByID(ctx context.Context, id exercise.ExerciseID, lang i18n.Language) (exercise.Exercise, error) {
	var model ExerciseModel
	if err := r.db.WithContext(ctx).Preload("Translations").Preload("Muscles.Translations").Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exercise.ErrExerciseNotFound
		}
		return nil, err
	}
	return model.ToDomain(lang), nil
}

func (r *exerciseRepository) Exists(ctx context.Context, id exercise.ExerciseID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&ExerciseModel{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
