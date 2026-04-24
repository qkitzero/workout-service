package exercise

import "context"

type ExerciseRepository interface {
	FindAll(ctx context.Context) ([]Exercise, error)
	FindByID(ctx context.Context, id ExerciseID) (Exercise, error)
	Exists(ctx context.Context, id ExerciseID) (bool, error)
}
