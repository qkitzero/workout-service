package set

import (
	"context"
	"time"

	"github.com/qkitzero/workout-service/internal/application/auth"
	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/set"
	"github.com/qkitzero/workout-service/internal/domain/user"
)

type SetUsecase interface {
	CreateSet(ctx context.Context, exerciseID exercise.ExerciseID, rep set.Rep, weight set.Weight, trainedAt time.Time) (set.Set, error)
	ListSets(ctx context.Context) ([]set.Set, error)
}

type setUsecase struct {
	authService  auth.AuthService
	setRepo      set.SetRepository
	exerciseRepo exercise.ExerciseRepository
}

func NewSetUsecase(authService auth.AuthService, setRepo set.SetRepository, exerciseRepo exercise.ExerciseRepository) SetUsecase {
	return &setUsecase{authService: authService, setRepo: setRepo, exerciseRepo: exerciseRepo}
}

func (u *setUsecase) CreateSet(ctx context.Context, exerciseID exercise.ExerciseID, rep set.Rep, weight set.Weight, trainedAt time.Time) (set.Set, error) {
	userID, err := u.authService.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	newUserID, err := user.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	exists, err := u.exerciseRepo.Exists(ctx, exerciseID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, exercise.ErrExerciseNotFound
	}

	newSet := set.NewSet(set.NewSetID(), newUserID, exerciseID, rep, weight, trainedAt, time.Now())

	if err := u.setRepo.Create(ctx, newSet); err != nil {
		return nil, err
	}

	return newSet, nil
}

func (u *setUsecase) ListSets(ctx context.Context) ([]set.Set, error) {
	userID, err := u.authService.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	newUserID, err := user.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	sets, err := u.setRepo.FindByUserID(ctx, newUserID)
	if err != nil {
		return nil, err
	}

	return sets, nil
}
