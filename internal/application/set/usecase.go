package set

import (
	"context"
	"time"

	"github.com/qkitzero/workout-service/internal/application/auth"
	"github.com/qkitzero/workout-service/internal/domain/set"
	"github.com/qkitzero/workout-service/internal/domain/user"
)

type SetUsecase interface {
	CreateSet(ctx context.Context, exercise string, rep int32, weight float64, trainedAt time.Time) (set.Set, error)
	ListSets(ctx context.Context) ([]set.Set, error)
}

type setUsecase struct {
	authService auth.AuthService
	setRepo     set.SetRepository
}

func NewSetUsecase(authService auth.AuthService, setRepo set.SetRepository) SetUsecase {
	return &setUsecase{authService: authService, setRepo: setRepo}
}

func (u *setUsecase) CreateSet(ctx context.Context, exercise string, rep int32, weight float64, trainedAt time.Time) (set.Set, error) {
	userID, err := u.authService.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	newUserID, err := user.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	newExercise, err := set.NewExercise(exercise)
	if err != nil {
		return nil, err
	}

	newRep, err := set.NewRep(rep)
	if err != nil {
		return nil, err
	}

	newWeight, err := set.NewWeight(weight)
	if err != nil {
		return nil, err
	}

	newSet := set.NewSet(set.NewSetID(), newUserID, newExercise, newRep, newWeight, trainedAt, time.Now())

	if err := u.setRepo.Create(newSet); err != nil {
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

	sets, err := u.setRepo.FindByUserID(newUserID)
	if err != nil {
		return nil, err
	}

	return sets, nil
}
