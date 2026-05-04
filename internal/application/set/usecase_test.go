package set

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/qkitzero/workout-service/internal/application/paging"
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
		name           string
		success        bool
		ctx            context.Context
		userID         string
		getUserErr     error
		setOwner       string
		findByIDErr    error
		workoutFinish  bool
		findWorkoutErr error
		existsResult   bool
		existsErr      error
		updateErr      error
	}{
		{"success update set", true, context.Background(), validUserID, nil, validUserID, nil, false, nil, true, nil, nil},
		{"failure get user error", false, context.Background(), "", fmt.Errorf("get user error"), "", nil, false, nil, false, nil, nil},
		{"failure empty user id", false, context.Background(), "", nil, "", nil, false, nil, false, nil, nil},
		{"failure set not found", false, context.Background(), validUserID, nil, "", set.ErrSetNotFound, false, nil, false, nil, nil},
		{"failure set forbidden", false, context.Background(), validUserID, nil, otherUserID, nil, false, nil, false, nil, nil},
		{"failure workout not found", false, context.Background(), validUserID, nil, validUserID, nil, false, workout.ErrWorkoutNotFound, false, nil, nil},
		{"failure workout already finished", false, context.Background(), validUserID, nil, validUserID, nil, true, nil, false, nil, nil},
		{"failure exercise not found", false, context.Background(), validUserID, nil, validUserID, nil, false, nil, false, nil, nil},
		{"failure exists error", false, context.Background(), validUserID, nil, validUserID, nil, false, nil, false, errors.New("exists error"), nil},
		{"failure update error", false, context.Background(), validUserID, nil, validUserID, nil, false, nil, true, nil, errors.New("update error")},
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

			mockWorkout := mocksworkout.NewMockWorkout(ctrl)
			mockWorkout.EXPECT().IsFinished().Return(tt.workoutFinish).AnyTimes()
			mockWorkoutRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(mockWorkout, tt.findWorkoutErr).AnyTimes()

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
		name           string
		success        bool
		ctx            context.Context
		userID         string
		getUserErr     error
		setOwner       string
		findByIDErr    error
		workoutFinish  bool
		findWorkoutErr error
		deleteErr      error
	}{
		{"success delete set", true, context.Background(), validUserID, nil, validUserID, nil, false, nil, nil},
		{"failure get user error", false, context.Background(), "", fmt.Errorf("get user error"), "", nil, false, nil, nil},
		{"failure empty user id", false, context.Background(), "", nil, "", nil, false, nil, nil},
		{"failure set not found", false, context.Background(), validUserID, nil, "", set.ErrSetNotFound, false, nil, nil},
		{"failure set forbidden", false, context.Background(), validUserID, nil, otherUserID, nil, false, nil, nil},
		{"failure workout not found", false, context.Background(), validUserID, nil, validUserID, nil, false, workout.ErrWorkoutNotFound, nil},
		{"failure workout already finished", false, context.Background(), validUserID, nil, validUserID, nil, true, nil, nil},
		{"failure delete error", false, context.Background(), validUserID, nil, validUserID, nil, false, nil, errors.New("delete error")},
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
			mockSet.EXPECT().WorkoutID().Return(workout.NewWorkoutID()).AnyTimes()
			mockSetRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(mockSet, tt.findByIDErr).AnyTimes()
			mockSetRepository.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(tt.deleteErr).AnyTimes()

			mockWorkout := mocksworkout.NewMockWorkout(ctrl)
			mockWorkout.EXPECT().IsFinished().Return(tt.workoutFinish).AnyTimes()
			mockWorkoutRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(mockWorkout, tt.findWorkoutErr).AnyTimes()

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
	validUserID := "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"

	mockSetSample := func(ctrl *gomock.Controller, trainedAt time.Time) *mocksset.MockSet {
		m := mocksset.NewMockSet(ctrl)
		m.EXPECT().ID().Return(set.NewSetID()).AnyTimes()
		m.EXPECT().TrainedAt().Return(trainedAt).AnyTimes()
		return m
	}

	validPageToken := paging.EncodeCursor(struct {
		TrainedAt time.Time `json:"t"`
		SetID     set.SetID `json:"s"`
	}{TrainedAt: time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC), SetID: set.NewSetID()})

	tests := []struct {
		name            string
		success         bool
		ctx             context.Context
		userID          string
		getUserErr      error
		pageSize        int
		pageToken       string
		repoResult      func(ctrl *gomock.Controller) []set.Set
		findByUserIDErr error
		wantNextToken   bool
	}{
		{
			name:       "success list sets",
			success:    true,
			ctx:        context.Background(),
			userID:     validUserID,
			pageSize:   10,
			repoResult: func(ctrl *gomock.Controller) []set.Set { return []set.Set{} },
		},
		{
			name:       "success default page size",
			success:    true,
			ctx:        context.Background(),
			userID:     validUserID,
			pageSize:   0,
			repoResult: func(ctrl *gomock.Controller) []set.Set { return []set.Set{} },
		},
		{
			name:       "success page size clamped to max",
			success:    true,
			ctx:        context.Background(),
			userID:     validUserID,
			pageSize:   1000,
			repoResult: func(ctrl *gomock.Controller) []set.Set { return []set.Set{} },
		},
		{
			name:       "success with valid page token",
			success:    true,
			ctx:        context.Background(),
			userID:     validUserID,
			pageSize:   10,
			pageToken:  validPageToken,
			repoResult: func(ctrl *gomock.Controller) []set.Set { return []set.Set{} },
		},
		{
			name:     "success with next page token",
			success:  true,
			ctx:      context.Background(),
			userID:   validUserID,
			pageSize: 2,
			repoResult: func(ctrl *gomock.Controller) []set.Set {
				return []set.Set{
					mockSetSample(ctrl, time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)),
					mockSetSample(ctrl, time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)),
					mockSetSample(ctrl, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				}
			},
			wantNextToken: true,
		},
		{"failure get user error", false, context.Background(), "", fmt.Errorf("get user error"), 10, "", func(ctrl *gomock.Controller) []set.Set { return nil }, nil, false},
		{"failure empty user id", false, context.Background(), "", nil, 10, "", func(ctrl *gomock.Controller) []set.Set { return nil }, nil, false},
		{"failure find by user id error", false, context.Background(), validUserID, nil, 10, "", func(ctrl *gomock.Controller) []set.Set { return nil }, errors.New("find by user id error"), false},
		{"failure invalid page token", false, context.Background(), validUserID, nil, 10, "!!!not-base64!!!", func(ctrl *gomock.Controller) []set.Set { return nil }, nil, false},
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
			mockSetRepository.EXPECT().FindByUserID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.repoResult(ctrl), tt.findByUserIDErr).AnyTimes()

			u := NewSetUsecase(mockUserService, mockSetRepository, mockWorkoutRepository, mockExerciseRepository)

			_, nextToken, err := u.ListSets(tt.ctx, nil, nil, tt.pageSize, tt.pageToken)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if tt.success {
				if tt.wantNextToken && nextToken == "" {
					t.Errorf("expected next token, got empty")
				}
				if !tt.wantNextToken && nextToken != "" {
					t.Errorf("expected empty next token, got %q", nextToken)
				}
			}
		})
	}
}
