package workout

import (
	"context"
	"errors"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	setv1 "github.com/qkitzero/workout-service/gen/go/set/v1"
	workoutv1 "github.com/qkitzero/workout-service/gen/go/workout/v1"
	appworkout "github.com/qkitzero/workout-service/internal/application/workout"
	"github.com/qkitzero/workout-service/internal/domain/workout"
)

type WorkoutHandler struct {
	workoutv1.UnimplementedWorkoutServiceServer
	workoutUsecase appworkout.WorkoutUsecase
}

func NewWorkoutHandler(workoutUsecase appworkout.WorkoutUsecase) *WorkoutHandler {
	return &WorkoutHandler{
		workoutUsecase: workoutUsecase,
	}
}

func (h *WorkoutHandler) StartWorkout(ctx context.Context, req *workoutv1.StartWorkoutRequest) (*workoutv1.StartWorkoutResponse, error) {
	w, err := h.workoutUsecase.StartWorkout(ctx)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		log.Printf("StartWorkout: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &workoutv1.StartWorkoutResponse{
		WorkoutId: w.ID().String(),
	}, nil
}

func (h *WorkoutHandler) FinishWorkout(ctx context.Context, req *workoutv1.FinishWorkoutRequest) (*workoutv1.FinishWorkoutResponse, error) {
	id, err := workout.NewWorkoutIDFromString(req.GetWorkoutId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	w, err := h.workoutUsecase.FinishWorkout(ctx, id)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		switch {
		case errors.Is(err, workout.ErrWorkoutNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, workout.ErrWorkoutForbidden):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case errors.Is(err, workout.ErrWorkoutAlreadyFinished):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		log.Printf("FinishWorkout: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &workoutv1.FinishWorkoutResponse{
		FinishedAt: timestamppb.New(*w.FinishedAt()),
	}, nil
}

func (h *WorkoutHandler) GetWorkout(ctx context.Context, req *workoutv1.GetWorkoutRequest) (*workoutv1.GetWorkoutResponse, error) {
	id, err := workout.NewWorkoutIDFromString(req.GetWorkoutId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	w, sets, err := h.workoutUsecase.GetWorkout(ctx, id)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		switch {
		case errors.Is(err, workout.ErrWorkoutNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, workout.ErrWorkoutForbidden):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		log.Printf("GetWorkout: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	setMessages := make([]*setv1.Set, 0, len(sets))
	for _, s := range sets {
		setMessages = append(setMessages, &setv1.Set{
			SetId:      s.ID().String(),
			WorkoutId:  s.WorkoutID().String(),
			ExerciseId: s.ExerciseID().String(),
			Rep:        s.Rep().Int32(),
			Weight:     s.Weight().Float64(),
			TrainedAt:  timestamppb.New(s.TrainedAt()),
			CreatedAt:  timestamppb.New(s.CreatedAt()),
		})
	}

	var finishedAt *timestamppb.Timestamp
	if f := w.FinishedAt(); f != nil {
		finishedAt = timestamppb.New(*f)
	}
	return &workoutv1.GetWorkoutResponse{
		Workout: &workoutv1.Workout{
			WorkoutId:  w.ID().String(),
			StartedAt:  timestamppb.New(w.StartedAt()),
			FinishedAt: finishedAt,
			CreatedAt:  timestamppb.New(w.CreatedAt()),
		},
		Sets: setMessages,
	}, nil
}

func (h *WorkoutHandler) ListWorkouts(ctx context.Context, req *workoutv1.ListWorkoutsRequest) (*workoutv1.ListWorkoutsResponse, error) {
	var from, to *time.Time
	if req.GetFrom() != nil {
		t := req.GetFrom().AsTime()
		from = &t
	}
	if req.GetTo() != nil {
		t := req.GetTo().AsTime()
		to = &t
	}
	if from != nil && to != nil && !from.Before(*to) {
		return nil, status.Error(codes.InvalidArgument, "from must be earlier than to")
	}

	workouts, err := h.workoutUsecase.ListWorkouts(ctx, from, to)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		log.Printf("ListWorkouts: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	responses := make([]*workoutv1.Workout, 0, len(workouts))
	for _, w := range workouts {
		var finishedAt *timestamppb.Timestamp
		if f := w.FinishedAt(); f != nil {
			finishedAt = timestamppb.New(*f)
		}
		responses = append(responses, &workoutv1.Workout{
			WorkoutId:  w.ID().String(),
			StartedAt:  timestamppb.New(w.StartedAt()),
			FinishedAt: finishedAt,
			CreatedAt:  timestamppb.New(w.CreatedAt()),
		})
	}

	return &workoutv1.ListWorkoutsResponse{
		Workouts: responses,
	}, nil
}
