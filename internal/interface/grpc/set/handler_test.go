package set

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	setv1 "github.com/qkitzero/workout-service/gen/go/set/v1"
	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/set"
	"github.com/qkitzero/workout-service/internal/domain/workout"
	mocksappset "github.com/qkitzero/workout-service/mocks/application/set"
	mocksset "github.com/qkitzero/workout-service/mocks/domain/set"
)

func TestCreateSet(t *testing.T) {
	t.Parallel()
	validWorkoutID := "a1a1a1a1-bbac-4bb9-a41d-b04f5afc4425"
	validExerciseID := "f1f538e5-4a37-409c-be99-09ee7bfefc50"
	trainedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		ctx          context.Context
		workoutID    string
		exerciseID   string
		rep          int32
		weight       float64
		callUsecase  bool
		createSetErr error
		wantCode     codes.Code
	}{
		{"success create set", context.Background(), validWorkoutID, validExerciseID, 10, 60.0, true, nil, codes.OK},
		{"failure invalid workout id", context.Background(), "not-a-uuid", validExerciseID, 10, 60.0, false, nil, codes.InvalidArgument},
		{"failure invalid exercise id", context.Background(), validWorkoutID, "not-a-uuid", 10, 60.0, false, nil, codes.InvalidArgument},
		{"failure invalid rep", context.Background(), validWorkoutID, validExerciseID, 0, 60.0, false, nil, codes.InvalidArgument},
		{"failure negative weight", context.Background(), validWorkoutID, validExerciseID, 10, -1.0, false, nil, codes.InvalidArgument},
		{"failure workout not found", context.Background(), validWorkoutID, validExerciseID, 10, 60.0, true, workout.ErrWorkoutNotFound, codes.NotFound},
		{"failure workout forbidden", context.Background(), validWorkoutID, validExerciseID, 10, 60.0, true, workout.ErrWorkoutForbidden, codes.PermissionDenied},
		{"failure workout already finished", context.Background(), validWorkoutID, validExerciseID, 10, 60.0, true, workout.ErrWorkoutAlreadyFinished, codes.FailedPrecondition},
		{"failure exercise not found", context.Background(), validWorkoutID, validExerciseID, 10, 60.0, true, exercise.ErrExerciseNotFound, codes.NotFound},
		{"failure usecase error", context.Background(), validWorkoutID, validExerciseID, 10, 60.0, true, fmt.Errorf("create set error"), codes.Internal},
		{"failure status preserved", context.Background(), validWorkoutID, validExerciseID, 10, 60.0, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappset.NewMockSetUsecase(ctrl)
			mockSet := mocksset.NewMockSet(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().CreateSet(tt.ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), trainedAt).Return(mockSet, tt.createSetErr).Times(1)
				mockSet.EXPECT().ID().Return(set.NewSetID()).AnyTimes()
			}

			handler := NewSetHandler(mockUsecase)

			req := &setv1.CreateSetRequest{
				WorkoutId:  tt.workoutID,
				ExerciseId: tt.exerciseID,
				Rep:        tt.rep,
				Weight:     tt.weight,
				TrainedAt:  timestamppb.New(trainedAt),
			}

			_, err := handler.CreateSet(tt.ctx, req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestGetSet(t *testing.T) {
	t.Parallel()
	validSetID := "f1f538e5-4a37-409c-be99-09ee7bfefc50"

	mockSetSample := func(ctrl *gomock.Controller) *mocksset.MockSet {
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

	tests := []struct {
		name        string
		ctx         context.Context
		setID       string
		callUsecase bool
		getSetErr   error
		wantCode    codes.Code
	}{
		{"success get set", context.Background(), validSetID, true, nil, codes.OK},
		{"failure invalid set id", context.Background(), "not-a-uuid", false, nil, codes.InvalidArgument},
		{"failure set not found", context.Background(), validSetID, true, set.ErrSetNotFound, codes.NotFound},
		{"failure set forbidden", context.Background(), validSetID, true, set.ErrSetForbidden, codes.PermissionDenied},
		{"failure usecase error", context.Background(), validSetID, true, fmt.Errorf("get set error"), codes.Internal},
		{"failure status preserved", context.Background(), validSetID, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappset.NewMockSetUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().GetSet(tt.ctx, gomock.Any()).Return(mockSetSample(ctrl), tt.getSetErr).Times(1)
			}

			handler := NewSetHandler(mockUsecase)

			_, err := handler.GetSet(tt.ctx, &setv1.GetSetRequest{SetId: tt.setID})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestUpdateSet(t *testing.T) {
	t.Parallel()
	validSetID := "f1f538e5-4a37-409c-be99-09ee7bfefc50"
	validExerciseID := "a1a1a1a1-bbac-4bb9-a41d-b04f5afc4425"
	trainedAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	mockSetSample := func(ctrl *gomock.Controller) *mocksset.MockSet {
		m := mocksset.NewMockSet(ctrl)
		m.EXPECT().ID().Return(set.NewSetID()).AnyTimes()
		m.EXPECT().WorkoutID().Return(workout.NewWorkoutID()).AnyTimes()
		m.EXPECT().ExerciseID().Return(exercise.NewExerciseID()).AnyTimes()
		m.EXPECT().Rep().Return(set.Rep(10)).AnyTimes()
		m.EXPECT().Weight().Return(set.Weight(60.0)).AnyTimes()
		m.EXPECT().TrainedAt().Return(trainedAt).AnyTimes()
		m.EXPECT().CreatedAt().Return(time.Now()).AnyTimes()
		return m
	}

	tests := []struct {
		name         string
		ctx          context.Context
		setID        string
		exerciseID   string
		rep          int32
		weight       float64
		callUsecase  bool
		updateSetErr error
		wantCode     codes.Code
	}{
		{"success update set", context.Background(), validSetID, validExerciseID, 10, 60.0, true, nil, codes.OK},
		{"failure invalid set id", context.Background(), "not-a-uuid", validExerciseID, 10, 60.0, false, nil, codes.InvalidArgument},
		{"failure invalid exercise id", context.Background(), validSetID, "not-a-uuid", 10, 60.0, false, nil, codes.InvalidArgument},
		{"failure invalid rep", context.Background(), validSetID, validExerciseID, 0, 60.0, false, nil, codes.InvalidArgument},
		{"failure negative weight", context.Background(), validSetID, validExerciseID, 10, -1.0, false, nil, codes.InvalidArgument},
		{"failure set not found", context.Background(), validSetID, validExerciseID, 10, 60.0, true, set.ErrSetNotFound, codes.NotFound},
		{"failure set forbidden", context.Background(), validSetID, validExerciseID, 10, 60.0, true, set.ErrSetForbidden, codes.PermissionDenied},
		{"failure workout not found", context.Background(), validSetID, validExerciseID, 10, 60.0, true, workout.ErrWorkoutNotFound, codes.NotFound},
		{"failure workout already finished", context.Background(), validSetID, validExerciseID, 10, 60.0, true, workout.ErrWorkoutAlreadyFinished, codes.FailedPrecondition},
		{"failure exercise not found", context.Background(), validSetID, validExerciseID, 10, 60.0, true, exercise.ErrExerciseNotFound, codes.NotFound},
		{"failure usecase error", context.Background(), validSetID, validExerciseID, 10, 60.0, true, fmt.Errorf("update set error"), codes.Internal},
		{"failure status preserved", context.Background(), validSetID, validExerciseID, 10, 60.0, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappset.NewMockSetUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().UpdateSet(tt.ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), trainedAt).Return(mockSetSample(ctrl), tt.updateSetErr).Times(1)
			}

			handler := NewSetHandler(mockUsecase)

			req := &setv1.UpdateSetRequest{
				SetId:      tt.setID,
				ExerciseId: tt.exerciseID,
				Rep:        tt.rep,
				Weight:     tt.weight,
				TrainedAt:  timestamppb.New(trainedAt),
			}

			_, err := handler.UpdateSet(tt.ctx, req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestDeleteSet(t *testing.T) {
	t.Parallel()
	validSetID := "f1f538e5-4a37-409c-be99-09ee7bfefc50"

	tests := []struct {
		name         string
		ctx          context.Context
		setID        string
		callUsecase  bool
		deleteSetErr error
		wantCode     codes.Code
	}{
		{"success delete set", context.Background(), validSetID, true, nil, codes.OK},
		{"failure invalid set id", context.Background(), "not-a-uuid", false, nil, codes.InvalidArgument},
		{"failure set not found", context.Background(), validSetID, true, set.ErrSetNotFound, codes.NotFound},
		{"failure set forbidden", context.Background(), validSetID, true, set.ErrSetForbidden, codes.PermissionDenied},
		{"failure workout not found", context.Background(), validSetID, true, workout.ErrWorkoutNotFound, codes.NotFound},
		{"failure workout already finished", context.Background(), validSetID, true, workout.ErrWorkoutAlreadyFinished, codes.FailedPrecondition},
		{"failure usecase error", context.Background(), validSetID, true, fmt.Errorf("delete set error"), codes.Internal},
		{"failure status preserved", context.Background(), validSetID, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappset.NewMockSetUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().DeleteSet(tt.ctx, gomock.Any()).Return(tt.deleteSetErr).Times(1)
			}

			handler := NewSetHandler(mockUsecase)

			_, err := handler.DeleteSet(tt.ctx, &setv1.DeleteSetRequest{SetId: tt.setID})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestListSets(t *testing.T) {
	t.Parallel()

	mockSetSample := func(ctrl *gomock.Controller) *mocksset.MockSet {
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

	tests := []struct {
		name        string
		ctx         context.Context
		listSetsErr error
		wantCode    codes.Code
	}{
		{"success list sets", context.Background(), nil, codes.OK},
		{"failure list sets error", context.Background(), fmt.Errorf("list sets error"), codes.Internal},
		{"failure status preserved", context.Background(), status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappset.NewMockSetUsecase(ctrl)
			mockUsecase.EXPECT().ListSets(tt.ctx).Return([]set.Set{mockSetSample(ctrl)}, tt.listSetsErr).Times(1)

			handler := NewSetHandler(mockUsecase)

			_, err := handler.ListSets(tt.ctx, &setv1.ListSetsRequest{})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}
