package set

import (
	"gorm.io/gorm"

	"github.com/qkitzero/workout-service/internal/domain/set"
)

type setRepository struct {
	db *gorm.DB
}

func NewSetRepository(db *gorm.DB) set.SetRepository {
	return &setRepository{db: db}
}

func (r *setRepository) Create(s set.Set) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		setModel := SetModel{
			ID:          s.ID(),
			UserID:      s.UserID(),
			Exercise:    s.Exercise(),
			Rep:         s.Rep(),
			Weight:      s.Weight(),
			TrainedAt:   s.TrainedAt(),
			CreatedAt:   s.CreatedAt(),
		}

		if err := tx.Create(&setModel).Error; err != nil {
			return err
		}

		return nil
	})
}
