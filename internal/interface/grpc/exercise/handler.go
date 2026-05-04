package exercise

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	exercisev1 "github.com/qkitzero/workout-service/gen/go/exercise/v1"
	appexercise "github.com/qkitzero/workout-service/internal/application/exercise"
	"github.com/qkitzero/workout-service/internal/domain/i18n"
)

type ExerciseHandler struct {
	exercisev1.UnimplementedExerciseServiceServer
	exerciseUsecase appexercise.ExerciseUsecase
}

func NewExerciseHandler(exerciseUsecase appexercise.ExerciseUsecase) *ExerciseHandler {
	return &ExerciseHandler{
		exerciseUsecase: exerciseUsecase,
	}
}

func (h *ExerciseHandler) ListExercises(ctx context.Context, req *exercisev1.ListExercisesRequest) (*exercisev1.ListExercisesResponse, error) {
	lang := i18n.LanguageJa
	if req.GetLang() != "" {
		parsed, err := i18n.NewLanguage(req.GetLang())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		lang = parsed
	}

	exercises, err := h.exerciseUsecase.ListExercises(ctx, lang)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		log.Printf("ListExercises: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	responses := make([]*exercisev1.Exercise, 0, len(exercises))
	for _, e := range exercises {
		domainMuscles := e.Muscles()
		muscles := make([]*exercisev1.Muscle, 0, len(domainMuscles))
		for _, m := range domainMuscles {
			muscles = append(muscles, &exercisev1.Muscle{
				MuscleId: m.ID().String(),
				Code:     m.Code().String(),
				Name:     m.Name().String(),
			})
		}
		responses = append(responses, &exercisev1.Exercise{
			ExerciseId: e.ID().String(),
			Code:       e.Code().String(),
			Name:       e.Name().String(),
			Category:   e.Category().String(),
			Muscles:    muscles,
		})
	}

	return &exercisev1.ListExercisesResponse{
		Exercises: responses,
	}, nil
}
