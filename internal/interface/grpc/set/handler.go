package set

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	setv1 "github.com/qkitzero/workout-service/gen/go/set/v1"
	appset "github.com/qkitzero/workout-service/internal/application/set"
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
	s, err := h.setUsecase.CreateSet(
		ctx,
		req.GetExercise(),
		req.GetRep(),
		req.GetWeight(),
		req.GetTrainedAt().AsTime(),
	)
	if err != nil {
		return nil, err
	}

	return &setv1.CreateSetResponse{
		SetId: s.ID().String(),
	}, nil
}

func (h *SetHandler) ListSets(ctx context.Context, req *setv1.ListSetsRequest) (*setv1.ListSetsResponse, error) {
	sets, err := h.setUsecase.ListSets(ctx)
	if err != nil {
		return nil, err
	}

	var responses []*setv1.Set
	for _, s := range sets {
		responses = append(responses, &setv1.Set{
			SetId:     s.ID().String(),
			Exercise:  s.Exercise().String(),
			Rep:       s.Rep().Int32(),
			Weight:    s.Weight().Float64(),
			TrainedAt: timestamppb.New(s.TrainedAt()),
			CreatedAt: timestamppb.New(s.CreatedAt()),
		})
	}

	return &setv1.ListSetsResponse{
		Sets: responses,
	}, nil
}
