package workout

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	workoutv1 "github.com/qkitzero/workout-service/gen/go/workout/v1"
	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/set"
	"github.com/qkitzero/workout-service/internal/domain/workout"
	mocksappworkout "github.com/qkitzero/workout-service/mocks/application/workout"
	mocksset "github.com/qkitzero/workout-service/mocks/domain/set"
	mocksworkout "github.com/qkitzero/workout-service/mocks/domain/workout"
)

func mockSetSample(ctrl *gomock.Controller) *mocksset.MockSet {
	m := mocksset.NewMockSet(ctrl)
	m.EXPECT().ID().Return(set.NewSetID()).AnyTimes()
	m.EXPECT().WorkoutID().Return(workout.NewWorkoutID()).AnyTimes()
	m.EXPECT().ExerciseID().Return(exercise.NewExerciseID()).AnyTimes()
	m.EXPECT().Rep().Return(set.Rep(10)).AnyTimes()
	m.EXPECT().Weight().Return(set.Weight(60.0)).AnyTimes()
	m.EXPECT().TrainedAt().Return(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)).AnyTimes()
	m.EXPECT().CreatedAt().Return(time.Now()).AnyTimes()
	return m
}

func mockWorkoutSample(ctrl *gomock.Controller) *mocksworkout.MockWorkout {
	m := mocksworkout.NewMockWorkout(ctrl)
	m.EXPECT().ID().Return(workout.NewWorkoutID()).AnyTimes()
	m.EXPECT().StartedAt().Return(time.Now()).AnyTimes()
	finishedAt := time.Now()
	m.EXPECT().FinishedAt().Return(&finishedAt).AnyTimes()
	m.EXPECT().CreatedAt().Return(time.Now()).AnyTimes()
	return m
}

func mockInProgressWorkoutSample(ctrl *gomock.Controller) *mocksworkout.MockWorkout {
	m := mocksworkout.NewMockWorkout(ctrl)
	m.EXPECT().ID().Return(workout.NewWorkoutID()).AnyTimes()
	m.EXPECT().StartedAt().Return(time.Now()).AnyTimes()
	m.EXPECT().FinishedAt().Return(nil).AnyTimes()
	m.EXPECT().CreatedAt().Return(time.Now()).AnyTimes()
	return m
}

func TestStartWorkout(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		ctx         context.Context
		callUsecase bool
		usecaseErr  error
		wantCode    codes.Code
	}{
		{"success start workout", context.Background(), true, nil, codes.OK},
		{"failure usecase error", context.Background(), true, fmt.Errorf("create error"), codes.Internal},
		{"failure status preserved", context.Background(), true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappworkout.NewMockWorkoutUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().StartWorkout(tt.ctx).Return(mockWorkoutSample(ctrl), tt.usecaseErr).Times(1)
			}

			handler := NewWorkoutHandler(mockUsecase)

			_, err := handler.StartWorkout(tt.ctx, &workoutv1.StartWorkoutRequest{})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestFinishWorkout(t *testing.T) {
	t.Parallel()
	validID := "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"

	tests := []struct {
		name        string
		ctx         context.Context
		workoutID   string
		callUsecase bool
		usecaseErr  error
		wantCode    codes.Code
	}{
		{"success finish workout", context.Background(), validID, true, nil, codes.OK},
		{"failure invalid workout id", context.Background(), "not-a-uuid", false, nil, codes.InvalidArgument},
		{"failure not found", context.Background(), validID, true, workout.ErrWorkoutNotFound, codes.NotFound},
		{"failure forbidden", context.Background(), validID, true, workout.ErrWorkoutForbidden, codes.PermissionDenied},
		{"failure already finished", context.Background(), validID, true, workout.ErrWorkoutAlreadyFinished, codes.FailedPrecondition},
		{"failure usecase error", context.Background(), validID, true, fmt.Errorf("update error"), codes.Internal},
		{"failure status preserved", context.Background(), validID, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappworkout.NewMockWorkoutUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().FinishWorkout(tt.ctx, gomock.Any()).Return(mockWorkoutSample(ctrl), tt.usecaseErr).Times(1)
			}

			handler := NewWorkoutHandler(mockUsecase)

			_, err := handler.FinishWorkout(tt.ctx, &workoutv1.FinishWorkoutRequest{WorkoutId: tt.workoutID})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestGetWorkout(t *testing.T) {
	t.Parallel()
	validID := "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"

	tests := []struct {
		name        string
		ctx         context.Context
		workoutID   string
		callUsecase bool
		usecaseErr  error
		wantCode    codes.Code
	}{
		{"success get workout", context.Background(), validID, true, nil, codes.OK},
		{"success in-progress workout", context.Background(), validID, true, nil, codes.OK},
		{"failure invalid workout id", context.Background(), "not-a-uuid", false, nil, codes.InvalidArgument},
		{"failure not found", context.Background(), validID, true, workout.ErrWorkoutNotFound, codes.NotFound},
		{"failure forbidden", context.Background(), validID, true, workout.ErrWorkoutForbidden, codes.PermissionDenied},
		{"failure usecase error", context.Background(), validID, true, fmt.Errorf("get error"), codes.Internal},
		{"failure status preserved", context.Background(), validID, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappworkout.NewMockWorkoutUsecase(ctrl)
			if tt.callUsecase {
				var w *mocksworkout.MockWorkout
				if tt.name == "success in-progress workout" {
					w = mockInProgressWorkoutSample(ctrl)
				} else {
					w = mockWorkoutSample(ctrl)
				}
				mockUsecase.EXPECT().GetWorkout(tt.ctx, gomock.Any()).Return(w, []set.Set{mockSetSample(ctrl)}, tt.usecaseErr).Times(1)
			}

			handler := NewWorkoutHandler(mockUsecase)

			_, err := handler.GetWorkout(tt.ctx, &workoutv1.GetWorkoutRequest{WorkoutId: tt.workoutID})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestListWorkouts(t *testing.T) {
	t.Parallel()
	validFrom := timestamppb.New(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	validTo := timestamppb.New(time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC))
	invalidTo := timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))

	tests := []struct {
		name        string
		ctx         context.Context
		req         *workoutv1.ListWorkoutsRequest
		callUsecase bool
		usecaseErr  error
		wantCode    codes.Code
	}{
		{"success no filters", context.Background(), &workoutv1.ListWorkoutsRequest{}, true, nil, codes.OK},
		{"success with filters", context.Background(), &workoutv1.ListWorkoutsRequest{From: validFrom, To: validTo}, true, nil, codes.OK},
		{"failure invalid from/to order", context.Background(), &workoutv1.ListWorkoutsRequest{From: validFrom, To: invalidTo}, false, nil, codes.InvalidArgument},
		{"failure usecase error", context.Background(), &workoutv1.ListWorkoutsRequest{}, true, fmt.Errorf("list error"), codes.Internal},
		{"failure status preserved", context.Background(), &workoutv1.ListWorkoutsRequest{}, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappworkout.NewMockWorkoutUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().ListWorkouts(tt.ctx, gomock.Any(), gomock.Any()).Return([]workout.Workout{mockWorkoutSample(ctrl)}, tt.usecaseErr).Times(1)
			}

			handler := NewWorkoutHandler(mockUsecase)

			_, err := handler.ListWorkouts(tt.ctx, tt.req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}
