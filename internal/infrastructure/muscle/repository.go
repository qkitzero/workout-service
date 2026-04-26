package muscle

import (
	"context"

	"gorm.io/gorm"

	"github.com/qkitzero/workout-service/internal/domain/i18n"
	"github.com/qkitzero/workout-service/internal/domain/muscle"
)

type muscleRepository struct {
	db *gorm.DB
}

func NewMuscleRepository(db *gorm.DB) muscle.MuscleRepository {
	return &muscleRepository{db: db}
}

func (r *muscleRepository) FindAll(ctx context.Context, lang i18n.Language) ([]muscle.Muscle, error) {
	var models []MuscleModel
	if err := r.db.WithContext(ctx).Preload("Translations").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]muscle.Muscle, len(models))
	for i, m := range models {
		result[i] = m.ToDomain(lang)
	}
	return result, nil
}
