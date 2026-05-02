package set

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	setv1 "github.com/qkitzero/workout-service/gen/go/set/v1"
	appset "github.com/qkitzero/workout-service/internal/application/set"
	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/set"
	"github.com/qkitzero/workout-service/internal/domain/workout"
)

type SetHandler struct {
	setv1.UnimplementedSetServiceServer
	setUsecase appset.SetUsecase
}

func NewSetHandler(setUsecase appset.SetUsecase) *SetHandler {
	return &SetHandler{
		setUsecase: setUsecase,
	}
}

func (h *SetHandler) CreateSet(ctx context.Context, req *setv1.CreateSetRequest) (*setv1.CreateSetResponse, error) {
	workoutID, err := workout.NewWorkoutIDFromString(req.GetWorkoutId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	exerciseID, err := exercise.NewExerciseIDFromString(req.GetExerciseId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rep, err := set.NewRep(req.GetRep())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	weight, err := set.NewWeight(req.GetWeight())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	s, err := h.setUsecase.CreateSet(ctx, workoutID, exerciseID, rep, weight, req.GetTrainedAt().AsTime())
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
		case errors.Is(err, exercise.ErrExerciseNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Printf("CreateSet: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &setv1.CreateSetResponse{
		SetId: s.ID().String(),
	}, nil
}

func (h *SetHandler) ListSets(ctx context.Context, req *setv1.ListSetsRequest) (*setv1.ListSetsResponse, error) {
	sets, err := h.setUsecase.ListSets(ctx)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		log.Printf("ListSets: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	responses := make([]*setv1.Set, 0, len(sets))
	for _, s := range sets {
		responses = append(responses, &setv1.Set{
			SetId:      s.ID().String(),
			WorkoutId:  s.WorkoutID().String(),
			ExerciseId: s.ExerciseID().String(),
			Rep:        s.Rep().Int32(),
			Weight:     s.Weight().Float64(),
			TrainedAt:  timestamppb.New(s.TrainedAt()),
			CreatedAt:  timestamppb.New(s.CreatedAt()),
		})
	}

	return &setv1.ListSetsResponse{
		Sets: responses,
	}, nil
}

func (h *SetHandler) GetSet(ctx context.Context, req *setv1.GetSetRequest) (*setv1.GetSetResponse, error) {
	setID, err := set.NewSetIDFromString(req.GetSetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	s, err := h.setUsecase.GetSet(ctx, setID)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		switch {
		case errors.Is(err, set.ErrSetNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, set.ErrSetForbidden):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		log.Printf("GetSet: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &setv1.GetSetResponse{
		Set: &setv1.Set{
			SetId:      s.ID().String(),
			WorkoutId:  s.WorkoutID().String(),
			ExerciseId: s.ExerciseID().String(),
			Rep:        s.Rep().Int32(),
			Weight:     s.Weight().Float64(),
			TrainedAt:  timestamppb.New(s.TrainedAt()),
			CreatedAt:  timestamppb.New(s.CreatedAt()),
		},
	}, nil
}

func (h *SetHandler) UpdateSet(ctx context.Context, req *setv1.UpdateSetRequest) (*setv1.UpdateSetResponse, error) {
	setID, err := set.NewSetIDFromString(req.GetSetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	exerciseID, err := exercise.NewExerciseIDFromString(req.GetExerciseId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rep, err := set.NewRep(req.GetRep())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	weight, err := set.NewWeight(req.GetWeight())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	s, err := h.setUsecase.UpdateSet(ctx, setID, exerciseID, rep, weight, req.GetTrainedAt().AsTime())
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		switch {
		case errors.Is(err, set.ErrSetNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, set.ErrSetForbidden):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case errors.Is(err, workout.ErrWorkoutNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, workout.ErrWorkoutAlreadyFinished):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		case errors.Is(err, exercise.ErrExerciseNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Printf("UpdateSet: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &setv1.UpdateSetResponse{
		Set: &setv1.Set{
			SetId:      s.ID().String(),
			WorkoutId:  s.WorkoutID().String(),
			ExerciseId: s.ExerciseID().String(),
			Rep:        s.Rep().Int32(),
			Weight:     s.Weight().Float64(),
			TrainedAt:  timestamppb.New(s.TrainedAt()),
			CreatedAt:  timestamppb.New(s.CreatedAt()),
		},
	}, nil
}

func (h *SetHandler) DeleteSet(ctx context.Context, req *setv1.DeleteSetRequest) (*setv1.DeleteSetResponse, error) {
	setID, err := set.NewSetIDFromString(req.GetSetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.setUsecase.DeleteSet(ctx, setID); err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		switch {
		case errors.Is(err, set.ErrSetNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, set.ErrSetForbidden):
			return nil, status.Error(codes.PermissionDenied, err.Error())
		case errors.Is(err, workout.ErrWorkoutNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, workout.ErrWorkoutAlreadyFinished):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		log.Printf("DeleteSet: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &setv1.DeleteSetResponse{}, nil
}
