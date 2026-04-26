package muscle

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	musclev1 "github.com/qkitzero/workout-service/gen/go/muscle/v1"
	appmuscle "github.com/qkitzero/workout-service/internal/application/muscle"
	"github.com/qkitzero/workout-service/internal/domain/i18n"
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
	lang := i18n.LanguageJa
	if req.GetLang() != "" {
		parsed, err := i18n.NewLanguage(req.GetLang())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		lang = parsed
	}

	muscles, err := h.muscleUsecase.ListMuscles(ctx, lang)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		log.Printf("ListMuscles: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
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
