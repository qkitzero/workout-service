package set

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	setv1 "github.com/qkitzero/workout-service/gen/go/set/v1"
	appset "github.com/qkitzero/workout-service/internal/application/set"
	"github.com/qkitzero/workout-service/internal/domain/exercise"
	"github.com/qkitzero/workout-service/internal/domain/set"
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

	s, err := h.setUsecase.CreateSet(ctx, exerciseID, rep, weight, req.GetTrainedAt().AsTime())
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		if errors.Is(err, exercise.ErrExerciseNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
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
		return nil, status.Error(codes.Internal, err.Error())
	}

	responses := make([]*setv1.Set, 0, len(sets))
	for _, s := range sets {
		responses = append(responses, &setv1.Set{
			SetId:      s.ID().String(),
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
