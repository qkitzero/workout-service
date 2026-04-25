package muscle

import (
	"context"

	musclev1 "github.com/qkitzero/workout-service/gen/go/muscle/v1"
	appmuscle "github.com/qkitzero/workout-service/internal/application/muscle"
)

type MuscleHandler struct {
	musclev1.UnimplementedMuscleServiceServer
	muscleUsecase appmuscle.MuscleUsecase
}

func NewMuscleHandler(muscleUsecase appmuscle.MuscleUsecase) *MuscleHandler {
	return &MuscleHandler{
		muscleUsecase: muscleUsecase,
	}
}

func (h *MuscleHandler) ListMuscles(ctx context.Context, req *musclev1.ListMusclesRequest) (*musclev1.ListMusclesResponse, error) {
	muscles, err := h.muscleUsecase.ListMuscles(ctx, req.GetLang())
	if err != nil {
		return nil, err
	}

	responses := make([]*musclev1.Muscle, 0, len(muscles))
	for _, m := range muscles {
		responses = append(responses, &musclev1.Muscle{
			MuscleId: m.ID().String(),
			Code:     m.Code().String(),
			Name:     m.Name().String(),
		})
	}

	return &musclev1.ListMusclesResponse{
		Muscles: responses,
	}, nil
}
