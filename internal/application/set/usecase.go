package set

import (
	"context"
	"time"

	"github.com/qkitzero/workout-service/internal/application/user"
	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/set"
	domainuser "github.com/qkitzero/workout-service/internal/domain/user"
	"github.com/qkitzero/workout-service/internal/domain/workout"
)

type SetUsecase interface {
	CreateSet(ctx context.Context, workoutID workout.WorkoutID, exerciseID exercise.ExerciseID, rep set.Rep, weight set.Weight, trainedAt time.Time) (set.Set, error)
	ListSets(ctx context.Context) ([]set.Set, error)
	GetSet(ctx context.Context, id set.SetID) (set.Set, error)
	UpdateSet(ctx context.Context, id set.SetID, exerciseID exercise.ExerciseID, rep set.Rep, weight set.Weight, trainedAt time.Time) (set.Set, error)
	DeleteSet(ctx context.Context, id set.SetID) error
}

type setUsecase struct {
	userService  user.UserService
	setRepo      set.SetRepository
	workoutRepo  workout.WorkoutRepository
	exerciseRepo exercise.ExerciseRepository
}

func NewSetUsecase(userService user.UserService, setRepo set.SetRepository, workoutRepo workout.WorkoutRepository, exerciseRepo exercise.ExerciseRepository) SetUsecase {
	return &setUsecase{userService: userService, setRepo: setRepo, workoutRepo: workoutRepo, exerciseRepo: exerciseRepo}
}

func (u *setUsecase) CreateSet(ctx context.Context, workoutID workout.WorkoutID, exerciseID exercise.ExerciseID, rep set.Rep, weight set.Weight, trainedAt time.Time) (set.Set, error) {
	userID, err := u.userService.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	newUserID, err := domainuser.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	w, err := u.workoutRepo.FindByID(ctx, workoutID)
	if err != nil {
		return nil, err
	}
	if w.UserID() != newUserID {
		return nil, workout.ErrWorkoutForbidden
	}
	if w.IsFinished() {
		return nil, workout.ErrWorkoutAlreadyFinished
	}

	exists, err := u.exerciseRepo.Exists(ctx, exerciseID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, exercise.ErrExerciseNotFound
	}

	newSet := set.NewSet(set.NewSetID(), newUserID, workoutID, exerciseID, rep, weight, trainedAt, time.Now())

	if err := u.setRepo.Create(ctx, newSet); err != nil {
		return nil, err
	}

	return newSet, nil
}

func (u *setUsecase) ListSets(ctx context.Context) ([]set.Set, error) {
	userID, err := u.userService.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	newUserID, err := domainuser.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	sets, err := u.setRepo.FindByUserID(ctx, newUserID)
	if err != nil {
		return nil, err
	}

	return sets, nil
}

func (u *setUsecase) GetSet(ctx context.Context, id set.SetID) (set.Set, error) {
	userID, err := u.userService.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	newUserID, err := domainuser.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	s, err := u.setRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.UserID() != newUserID {
		return nil, set.ErrSetForbidden
	}

	return s, nil
}

func (u *setUsecase) UpdateSet(ctx context.Context, id set.SetID, exerciseID exercise.ExerciseID, rep set.Rep, weight set.Weight, trainedAt time.Time) (set.Set, error) {
	userID, err := u.userService.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	newUserID, err := domainuser.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	s, err := u.setRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.UserID() != newUserID {
		return nil, set.ErrSetForbidden
	}

	exists, err := u.exerciseRepo.Exists(ctx, exerciseID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, exercise.ErrExerciseNotFound
	}

	updated := set.NewSet(s.ID(), s.UserID(), s.WorkoutID(), exerciseID, rep, weight, trainedAt, s.CreatedAt())

	if err := u.setRepo.Update(ctx, updated); err != nil {
		return nil, err
	}

	return updated, nil
}

func (u *setUsecase) DeleteSet(ctx context.Context, id set.SetID) error {
	userID, err := u.userService.GetUser(ctx)
	if err != nil {
		return err
	}

	newUserID, err := domainuser.NewUserID(userID)
	if err != nil {
		return err
	}

	s, err := u.setRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if s.UserID() != newUserID {
		return set.ErrSetForbidden
	}

	return u.setRepo.Delete(ctx, id)
}
