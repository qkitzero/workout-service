package set

import (
	"context"

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
