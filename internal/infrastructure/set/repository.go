package set

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/qkitzero/workout-service/internal/domain/set"
	"github.com/qkitzero/workout-service/internal/domain/user"
	"github.com/qkitzero/workout-service/internal/domain/workout"
)

type setRepository struct {
	db *gorm.DB
}

func NewSetRepository(db *gorm.DB) set.SetRepository {
	return &setRepository{db: db}
}

func (r *setRepository) Create(ctx context.Context, s set.Set) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		setModel := SetModel{
			ID:         s.ID(),
			UserID:     s.UserID(),
			WorkoutID:  s.WorkoutID(),
			ExerciseID: s.ExerciseID(),
			Rep:        s.Rep(),
			Weight:     s.Weight(),
			TrainedAt:  s.TrainedAt(),
			CreatedAt:  s.CreatedAt(),
		}

		if err := tx.Create(&setModel).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *setRepository) Update(ctx context.Context, s set.Set) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		setModel := SetModel{
			ID:         s.ID(),
			UserID:     s.UserID(),
			WorkoutID:  s.WorkoutID(),
			ExerciseID: s.ExerciseID(),
			Rep:        s.Rep(),
			Weight:     s.Weight(),
			TrainedAt:  s.TrainedAt(),
			CreatedAt:  s.CreatedAt(),
		}

		if err := tx.Save(&setModel).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *setRepository) Delete(ctx context.Context, id set.SetID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&SetModel{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *setRepository) FindByID(ctx context.Context, id set.SetID) (set.Set, error) {
	var m SetModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, set.ErrSetNotFound
		}
		return nil, err
	}
	return set.NewSet(m.ID, m.UserID, m.WorkoutID, m.ExerciseID, m.Rep, m.Weight, m.TrainedAt, m.CreatedAt), nil
}

func (r *setRepository) FindByUserID(
	ctx context.Context,
	userID user.UserID,
	from, to *time.Time,
	limit int,
	cursorTrainedAt *time.Time,
	cursorSetID *set.SetID,
) ([]set.Set, error) {
	q := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if from != nil {
		q = q.Where("trained_at >= ?", *from)
	}
	if to != nil {
		q = q.Where("trained_at < ?", *to)
	}
	if cursorTrainedAt != nil && cursorSetID != nil {
		q = q.Where("(trained_at, id) < (?, ?)", *cursorTrainedAt, *cursorSetID)
	}

	q = q.Order("trained_at DESC").Order("id DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}

	var setModels []SetModel
	if err := q.Find(&setModels).Error; err != nil {
		return nil, err
	}

	sets := make([]set.Set, len(setModels))
	for i, m := range setModels {
		sets[i] = set.NewSet(m.ID, m.UserID, m.WorkoutID, m.ExerciseID, m.Rep, m.Weight, m.TrainedAt, m.CreatedAt)
	}

	return sets, nil
}

func (r *setRepository) FindByWorkoutID(ctx context.Context, workoutID workout.WorkoutID) ([]set.Set, error) {
	var setModels []SetModel
	if err := r.db.WithContext(ctx).Where("workout_id = ?", workoutID).Find(&setModels).Error; err != nil {
		return nil, err
	}

	sets := make([]set.Set, len(setModels))
	for i, m := range setModels {
		sets[i] = set.NewSet(m.ID, m.UserID, m.WorkoutID, m.ExerciseID, m.Rep, m.Weight, m.TrainedAt, m.CreatedAt)
	}

	return sets, nil
}
