package set

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/set"
	"github.com/qkitzero/workout-service/internal/domain/user"
	"github.com/qkitzero/workout-service/internal/domain/workout"
	mocksappuser "github.com/qkitzero/workout-service/mocks/application/user"
	mocksexercise "github.com/qkitzero/workout-service/mocks/domain/exercise"
	mocksset "github.com/qkitzero/workout-service/mocks/domain/set"
	mocksworkout "github.com/qkitzero/workout-service/mocks/domain/workout"
)

func TestCreateSet(t *testing.T) {
	t.Parallel()
	validUserID := "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"
	otherUserID := "11111111-bbac-4bb9-a41d-b04f5afc4425"
	workoutID := workout.NewWorkoutID()
	exerciseID := exercise.NewExerciseID()
	rep, _ := set.NewRep(10)
	weight, _ := set.NewWeight(60.0)
	trainedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		success       bool
		ctx           context.Context
		userID        string
		getUserErr    error
		workoutOwner  string
		workoutFinish bool
		findByIDErr   error
		existsResult  bool
		existsErr     error
		createErr     error
	}{
		{"success create set", true, context.Background(), validUserID, nil, validUserID, false, nil, true, nil, nil},
		{"failure get user error", false, context.Background(), "", fmt.Errorf("get user error"), "", false, nil, false, nil, nil},
		{"failure empty user id", false, context.Background(), "", nil, "", false, nil, false, nil, nil},
		{"failure workout not found", false, context.Background(), validUserID, nil, "", false, workout.ErrWorkoutNotFound, false, nil, nil},
		{"failure workout forbidden", false, context.Background(), validUserID, nil, otherUserID, false, nil, false, nil, nil},
		{"failure workout already finished", false, context.Background(), validUserID, nil, validUserID, true, nil, false, nil, nil},
		{"failure exercise not found", false, context.Background(), validUserID, nil, validUserID, false, nil, false, nil, nil},
		{"failure exists error", false, context.Background(), validUserID, nil, validUserID, false, nil, false, errors.New("exists error"), nil},
		{"failure create error", false, context.Background(), validUserID, nil, validUserID, false, nil, true, nil, errors.New("create error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockSetRepository := mocksset.NewMockSetRepository(ctrl)
			mockWorkoutRepository := mocksworkout.NewMockWorkoutRepository(ctrl)
			mockExerciseRepository := mocksexercise.NewMockExerciseRepository(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()

			mockWorkout := mocksworkout.NewMockWorkout(ctrl)
			mockWorkout.EXPECT().UserID().Return(user.UserID(tt.workoutOwner)).AnyTimes()
			mockWorkout.EXPECT().IsFinished().Return(tt.workoutFinish).AnyTimes()
			mockWorkoutRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(mockWorkout, tt.findByIDErr).AnyTimes()

			mockExerciseRepository.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(tt.existsResult, tt.existsErr).AnyTimes()
			mockSetRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tt.createErr).AnyTimes()

			u := NewSetUsecase(mockUserService, mockSetRepository, mockWorkoutRepository, mockExerciseRepository)

			_, err := u.CreateSet(tt.ctx, workoutID, exerciseID, rep, weight, trainedAt)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestGetSet(t *testing.T) {
	t.Parallel()
	validUserID := "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"
	otherUserID := "11111111-bbac-4bb9-a41d-b04f5afc4425"
	setID := set.NewSetID()

	tests := []struct {
		name        string
		success     bool
		ctx         context.Context
		userID      string
		getUserErr  error
		setOwner    string
		findByIDErr error
	}{
		{"success get set", true, context.Background(), validUserID, nil, validUserID, nil},
		{"failure get user error", false, context.Background(), "", fmt.Errorf("get user error"), "", nil},
		{"failure empty user id", false, context.Background(), "", nil, "", nil},
		{"failure set not found", false, context.Background(), validUserID, nil, "", set.ErrSetNotFound},
		{"failure set forbidden", false, context.Background(), validUserID, nil, otherUserID, nil},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockSetRepository := mocksset.NewMockSetRepository(ctrl)
			mockWorkoutRepository := mocksworkout.NewMockWorkoutRepository(ctrl)
			mockExerciseRepository := mocksexercise.NewMockExerciseRepository(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()

			mockSet := mocksset.NewMockSet(ctrl)
			mockSet.EXPECT().UserID().Return(user.UserID(tt.setOwner)).AnyTimes()
			mockSetRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(mockSet, tt.findByIDErr).AnyTimes()

			u := NewSetUsecase(mockUserService, mockSetRepository, mockWorkoutRepository, mockExerciseRepository)

			_, err := u.GetSet(tt.ctx, setID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestUpdateSet(t *testing.T) {
	t.Parallel()
	validUserID := "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"
	otherUserID := "11111111-bbac-4bb9-a41d-b04f5afc4425"
	setID := set.NewSetID()
	exerciseID := exercise.NewExerciseID()
	rep, _ := set.NewRep(10)
	weight, _ := set.NewWeight(60.0)
	trainedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		success      bool
		ctx          context.Context
		userID       string
		getUserErr   error
		setOwner     string
		findByIDErr  error
		existsResult bool
		existsErr    error
		updateErr    error
	}{
		{"success update set", true, context.Background(), validUserID, nil, validUserID, nil, true, nil, nil},
		{"failure get user error", false, context.Background(), "", fmt.Errorf("get user error"), "", nil, false, nil, nil},
		{"failure empty user id", false, context.Background(), "", nil, "", nil, false, nil, nil},
		{"failure set not found", false, context.Background(), validUserID, nil, "", set.ErrSetNotFound, false, nil, nil},
		{"failure set forbidden", false, context.Background(), validUserID, nil, otherUserID, nil, false, nil, nil},
		{"failure exercise not found", false, context.Background(), validUserID, nil, validUserID, nil, false, nil, nil},
		{"failure exists error", false, context.Background(), validUserID, nil, validUserID, nil, false, errors.New("exists error"), nil},
		{"failure update error", false, context.Background(), validUserID, nil, validUserID, nil, true, nil, errors.New("update error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockSetRepository := mocksset.NewMockSetRepository(ctrl)
			mockWorkoutRepository := mocksworkout.NewMockWorkoutRepository(ctrl)
			mockExerciseRepository := mocksexercise.NewMockExerciseRepository(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()

			mockSet := mocksset.NewMockSet(ctrl)
			mockSet.EXPECT().ID().Return(setID).AnyTimes()
			mockSet.EXPECT().UserID().Return(user.UserID(tt.setOwner)).AnyTimes()
			mockSet.EXPECT().WorkoutID().Return(workout.NewWorkoutID()).AnyTimes()
			mockSet.EXPECT().CreatedAt().Return(time.Now()).AnyTimes()
			mockSetRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(mockSet, tt.findByIDErr).AnyTimes()

			mockExerciseRepository.EXPECT().Exists(gomock.Any(), gomock.Any()).Return(tt.existsResult, tt.existsErr).AnyTimes()
			mockSetRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(tt.updateErr).AnyTimes()

			u := NewSetUsecase(mockUserService, mockSetRepository, mockWorkoutRepository, mockExerciseRepository)

			_, err := u.UpdateSet(tt.ctx, setID, exerciseID, rep, weight, trainedAt)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestDeleteSet(t *testing.T) {
	t.Parallel()
	validUserID := "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"
	otherUserID := "11111111-bbac-4bb9-a41d-b04f5afc4425"
	setID := set.NewSetID()

	tests := []struct {
		name        string
		success     bool
		ctx         context.Context
		userID      string
		getUserErr  error
		setOwner    string
		findByIDErr error
		deleteErr   error
	}{
		{"success delete set", true, context.Background(), validUserID, nil, validUserID, nil, nil},
		{"failure get user error", false, context.Background(), "", fmt.Errorf("get user error"), "", nil, nil},
		{"failure empty user id", false, context.Background(), "", nil, "", nil, nil},
		{"failure set not found", false, context.Background(), validUserID, nil, "", set.ErrSetNotFound, nil},
		{"failure set forbidden", false, context.Background(), validUserID, nil, otherUserID, nil, nil},
		{"failure delete error", false, context.Background(), validUserID, nil, validUserID, nil, errors.New("delete error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockSetRepository := mocksset.NewMockSetRepository(ctrl)
			mockWorkoutRepository := mocksworkout.NewMockWorkoutRepository(ctrl)
			mockExerciseRepository := mocksexercise.NewMockExerciseRepository(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()

			mockSet := mocksset.NewMockSet(ctrl)
			mockSet.EXPECT().UserID().Return(user.UserID(tt.setOwner)).AnyTimes()
			mockSetRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(mockSet, tt.findByIDErr).AnyTimes()
			mockSetRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(tt.deleteErr).AnyTimes()

			u := NewSetUsecase(mockUserService, mockSetRepository, mockWorkoutRepository, mockExerciseRepository)

			err := u.DeleteSet(tt.ctx, setID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestListSets(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		success         bool
		ctx             context.Context
		userID          string
		getUserErr      error
		findByUserIDErr error
	}{
		{"success list sets", true, context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", nil, nil},
		{"failure get user error", false, context.Background(), "", fmt.Errorf("get user error"), nil},
		{"failure empty user id", false, context.Background(), "", nil, nil},
		{"failure find by user id error", false, context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", nil, errors.New("find by user id error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockSetRepository := mocksset.NewMockSetRepository(ctrl)
			mockWorkoutRepository := mocksworkout.NewMockWorkoutRepository(ctrl)
			mockExerciseRepository := mocksexercise.NewMockExerciseRepository(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()
			mockSetRepository.EXPECT().FindByUserID(gomock.Any(), gomock.Any()).Return([]set.Set{}, tt.findByUserIDErr).AnyTimes()

			u := NewSetUsecase(mockUserService, mockSetRepository, mockWorkoutRepository, mockExerciseRepository)

			_, err := u.ListSets(tt.ctx)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
