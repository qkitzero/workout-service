package set

import (
	"gorm.io/gorm"

	"github.com/qkitzero/workout-service/internal/domain/set"
	"github.com/qkitzero/workout-service/internal/domain/user"
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

func (r *setRepository) FindByUserID(userID user.UserID) ([]set.Set, error) {
	var setModels []SetModel
	if err := r.db.Where("user_id = ?", userID).Find(&setModels).Error; err != nil {
		return nil, err
	}

	sets := make([]set.Set, len(setModels))
	for i, m := range setModels {
		sets[i] = set.NewSet(m.ID, m.UserID, m.Exercise, m.Rep, m.Weight, m.TrainedAt, m.CreatedAt)
	}

	return sets, nil
}
