package workout

import (
	"context"
	"time"

	"github.com/qkitzero/workout-service/internal/application/user"
	"github.com/qkitzero/workout-service/internal/domain/set"
	domainuser "github.com/qkitzero/workout-service/internal/domain/user"
	"github.com/qkitzero/workout-service/internal/domain/workout"
)

type WorkoutUsecase interface {
	StartWorkout(ctx context.Context) (workout.Workout, error)
	FinishWorkout(ctx context.Context, id workout.WorkoutID) (workout.Workout, error)
	GetWorkout(ctx context.Context, id workout.WorkoutID) (workout.Workout, []set.Set, error)
	ListWorkouts(ctx context.Context, from, to *time.Time) ([]workout.Workout, error)
}

type workoutUsecase struct {
	userService user.UserService
	workoutRepo workout.WorkoutRepository
	setRepo     set.SetRepository
}

func NewWorkoutUsecase(userService user.UserService, workoutRepo workout.WorkoutRepository, setRepo set.SetRepository) WorkoutUsecase {
	return &workoutUsecase{userService: userService, workoutRepo: workoutRepo, setRepo: setRepo}
}

func (u *workoutUsecase) StartWorkout(ctx context.Context) (workout.Workout, error) {
	userID, err := u.userService.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	newUserID, err := domainuser.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	w := workout.NewWorkout(workout.NewWorkoutID(), newUserID, now, nil, now)

	if err := u.workoutRepo.Create(ctx, w); err != nil {
		return nil, err
	}

	return w, nil
}

func (u *workoutUsecase) FinishWorkout(ctx context.Context, id workout.WorkoutID) (workout.Workout, error) {
	userID, err := u.userService.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	newUserID, err := domainuser.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	w, err := u.workoutRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if w.UserID() != newUserID {
		return nil, workout.ErrWorkoutForbidden
	}
	if w.IsFinished() {
		return nil, workout.ErrWorkoutAlreadyFinished
	}

	now := time.Now()
	finished := workout.NewWorkout(w.ID(), w.UserID(), w.StartedAt(), &now, w.CreatedAt())
	if err := u.workoutRepo.Update(ctx, finished); err != nil {
		return nil, err
	}

	return finished, nil
}

func (u *workoutUsecase) GetWorkout(ctx context.Context, id workout.WorkoutID) (workout.Workout, []set.Set, error) {
	userID, err := u.userService.GetUser(ctx)
	if err != nil {
		return nil, nil, err
	}

	newUserID, err := domainuser.NewUserID(userID)
	if err != nil {
		return nil, nil, err
	}

	w, err := u.workoutRepo.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if w.UserID() != newUserID {
		return nil, nil, workout.ErrWorkoutForbidden
	}

	sets, err := u.setRepo.FindByWorkoutID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	return w, sets, nil
}

func (u *workoutUsecase) ListWorkouts(ctx context.Context, from, to *time.Time) ([]workout.Workout, error) {
	userID, err := u.userService.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	newUserID, err := domainuser.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	return u.workoutRepo.FindByUserID(ctx, newUserID, from, to)
}
