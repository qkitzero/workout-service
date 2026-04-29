package workout

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/qkitzero/workout-service/internal/domain/user"
	"github.com/qkitzero/workout-service/internal/domain/workout"
)

type workoutRepository struct {
	db *gorm.DB
}

func NewWorkoutRepository(db *gorm.DB) workout.WorkoutRepository {
	return &workoutRepository{db: db}
}

func (r *workoutRepository) Create(ctx context.Context, w workout.Workout) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		model := WorkoutModel{
			ID:         w.ID(),
			UserID:     w.UserID(),
			StartedAt:  w.StartedAt(),
			FinishedAt: w.FinishedAt(),
			CreatedAt:  w.CreatedAt(),
		}

		if err := tx.Create(&model).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *workoutRepository) Update(ctx context.Context, w workout.Workout) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		model := WorkoutModel{
			ID:         w.ID(),
			UserID:     w.UserID(),
			StartedAt:  w.StartedAt(),
			FinishedAt: w.FinishedAt(),
			CreatedAt:  w.CreatedAt(),
		}

		if err := tx.Save(&model).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *workoutRepository) FindByID(ctx context.Context, id workout.WorkoutID) (workout.Workout, error) {
	var model WorkoutModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, workout.ErrWorkoutNotFound
		}
		return nil, err
	}
	return workout.NewWorkout(model.ID, model.UserID, model.StartedAt, model.FinishedAt, model.CreatedAt), nil
}

func (r *workoutRepository) FindByUserID(ctx context.Context, userID user.UserID, from, to *time.Time) ([]workout.Workout, error) {
	q := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if from != nil {
		q = q.Where("started_at >= ?", *from)
	}
	if to != nil {
		q = q.Where("started_at < ?", *to)
	}

	var models []WorkoutModel
	if err := q.Find(&models).Error; err != nil {
		return nil, err
	}

	workouts := make([]workout.Workout, len(models))
	for i, m := range models {
		workouts[i] = workout.NewWorkout(m.ID, m.UserID, m.StartedAt, m.FinishedAt, m.CreatedAt)
	}
	return workouts, nil
}

func (r *workoutRepository) Exists(ctx context.Context, id workout.WorkoutID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&WorkoutModel{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
