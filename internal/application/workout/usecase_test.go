package workout

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/qkitzero/workout-service/internal/domain/set"
	"github.com/qkitzero/workout-service/internal/domain/user"
	"github.com/qkitzero/workout-service/internal/domain/workout"
	mocksappuser "github.com/qkitzero/workout-service/mocks/application/user"
	mocksset "github.com/qkitzero/workout-service/mocks/domain/set"
	mocksworkout "github.com/qkitzero/workout-service/mocks/domain/workout"
)

func TestStartWorkout(t *testing.T) {
	t.Parallel()
	validUserID := "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"

	tests := []struct {
		name       string
		success    bool
		ctx        context.Context
		userID     string
		getUserErr error
		createErr  error
	}{
		{"success start workout", true, context.Background(), validUserID, nil, nil},
		{"failure get user error", false, context.Background(), "", fmt.Errorf("get user error"), nil},
		{"failure empty user id", false, context.Background(), "", nil, nil},
		{"failure create error", false, context.Background(), validUserID, nil, errors.New("create error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockWorkoutRepository := mocksworkout.NewMockWorkoutRepository(ctrl)
			mockSetRepository := mocksset.NewMockSetRepository(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()
			mockWorkoutRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tt.createErr).AnyTimes()

			u := NewWorkoutUsecase(mockUserService, mockWorkoutRepository, mockSetRepository)

			_, err := u.StartWorkout(tt.ctx)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestFinishWorkout(t *testing.T) {
	t.Parallel()
	validUserID := "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"
	otherUserID := "11111111-bbac-4bb9-a41d-b04f5afc4425"
	id := workout.NewWorkoutID()

	tests := []struct {
		name        string
		success     bool
		ctx         context.Context
		userID      string
		getUserErr  error
		owner       string
		finished    bool
		findByIDErr error
		updateErr   error
	}{
		{"success finish workout", true, context.Background(), validUserID, nil, validUserID, false, nil, nil},
		{"failure get user error", false, context.Background(), "", fmt.Errorf("get user error"), "", false, nil, nil},
		{"failure empty user id", false, context.Background(), "", nil, "", false, nil, nil},
		{"failure not found", false, context.Background(), validUserID, nil, "", false, workout.ErrWorkoutNotFound, nil},
		{"failure forbidden", false, context.Background(), validUserID, nil, otherUserID, false, nil, nil},
		{"failure already finished", false, context.Background(), validUserID, nil, validUserID, true, nil, nil},
		{"failure update error", false, context.Background(), validUserID, nil, validUserID, false, nil, errors.New("update error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockWorkoutRepository := mocksworkout.NewMockWorkoutRepository(ctrl)
			mockSetRepository := mocksset.NewMockSetRepository(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()

			mockWorkout := mocksworkout.NewMockWorkout(ctrl)
			mockWorkout.EXPECT().ID().Return(id).AnyTimes()
			mockWorkout.EXPECT().UserID().Return(user.UserID(tt.owner)).AnyTimes()
			mockWorkout.EXPECT().StartedAt().Return(time.Now()).AnyTimes()
			mockWorkout.EXPECT().CreatedAt().Return(time.Now()).AnyTimes()
			mockWorkout.EXPECT().IsFinished().Return(tt.finished).AnyTimes()
			mockWorkoutRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(mockWorkout, tt.findByIDErr).AnyTimes()
			mockWorkoutRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(tt.updateErr).AnyTimes()

			u := NewWorkoutUsecase(mockUserService, mockWorkoutRepository, mockSetRepository)

			_, err := u.FinishWorkout(tt.ctx, id)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestGetWorkout(t *testing.T) {
	t.Parallel()
	validUserID := "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"
	otherUserID := "11111111-bbac-4bb9-a41d-b04f5afc4425"
	id := workout.NewWorkoutID()

	tests := []struct {
		name             string
		success          bool
		ctx              context.Context
		userID           string
		getUserErr       error
		owner            string
		findByIDErr      error
		findByWorkoutErr error
	}{
		{"success get workout", true, context.Background(), validUserID, nil, validUserID, nil, nil},
		{"failure get user error", false, context.Background(), "", fmt.Errorf("get user error"), "", nil, nil},
		{"failure empty user id", false, context.Background(), "", nil, "", nil, nil},
		{"failure not found", false, context.Background(), validUserID, nil, "", workout.ErrWorkoutNotFound, nil},
		{"failure forbidden", false, context.Background(), validUserID, nil, otherUserID, nil, nil},
		{"failure find sets error", false, context.Background(), validUserID, nil, validUserID, nil, errors.New("sets error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockWorkoutRepository := mocksworkout.NewMockWorkoutRepository(ctrl)
			mockSetRepository := mocksset.NewMockSetRepository(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()

			mockWorkout := mocksworkout.NewMockWorkout(ctrl)
			mockWorkout.EXPECT().UserID().Return(user.UserID(tt.owner)).AnyTimes()
			mockWorkoutRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(mockWorkout, tt.findByIDErr).AnyTimes()
			mockSetRepository.EXPECT().FindByWorkoutID(gomock.Any(), gomock.Any()).Return([]set.Set{}, tt.findByWorkoutErr).AnyTimes()

			u := NewWorkoutUsecase(mockUserService, mockWorkoutRepository, mockSetRepository)

			_, _, err := u.GetWorkout(tt.ctx, id)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestListWorkouts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		success    bool
		ctx        context.Context
		userID     string
		getUserErr error
		findErr    error
	}{
		{"success list workouts", true, context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", nil, nil},
		{"failure get user error", false, context.Background(), "", fmt.Errorf("get user error"), nil},
		{"failure empty user id", false, context.Background(), "", nil, nil},
		{"failure find error", false, context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", nil, errors.New("find error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserService := mocksappuser.NewMockUserService(ctrl)
			mockWorkoutRepository := mocksworkout.NewMockWorkoutRepository(ctrl)
			mockSetRepository := mocksset.NewMockSetRepository(ctrl)
			mockUserService.EXPECT().GetUser(tt.ctx).Return(tt.userID, tt.getUserErr).AnyTimes()
			mockWorkoutRepository.EXPECT().FindByUserID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]workout.Workout{}, tt.findErr).AnyTimes()

			u := NewWorkoutUsecase(mockUserService, mockWorkoutRepository, mockSetRepository)

			_, err := u.ListWorkouts(tt.ctx, nil, nil)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
