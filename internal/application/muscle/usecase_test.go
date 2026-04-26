package muscle

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/qkitzero/workout-service/internal/domain/i18n"
	"github.com/qkitzero/workout-service/internal/domain/muscle"
	mocksmuscle "github.com/qkitzero/workout-service/mocks/domain/muscle"
)

func TestListMuscles(t *testing.T) {
	t.Parallel()

	code, err := muscle.NewCode("chest")
	if err != nil {
		t.Errorf("failed to new code: %v", err)
	}
	name, err := muscle.NewName("胸")
	if err != nil {
		t.Errorf("failed to new name: %v", err)
	}
	sample := muscle.NewMuscle(muscle.NewMuscleID(), code, name)

	tests := []struct {
		name        string
		success     bool
		ctx         context.Context
		lang        i18n.Language
		findAllResp []muscle.Muscle
		findAllErr  error
	}{
		{"success ja", true, context.Background(), i18n.LanguageJa, []muscle.Muscle{sample}, nil},
		{"success en", true, context.Background(), i18n.Language("en"), []muscle.Muscle{sample}, nil},
		{"success empty result", true, context.Background(), i18n.LanguageJa, []muscle.Muscle{}, nil},
		{"failure find all error", false, context.Background(), i18n.LanguageJa, nil, errors.New("find all error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockMuscleRepository := mocksmuscle.NewMockMuscleRepository(ctrl)
			mockMuscleRepository.EXPECT().FindAll(gomock.Any(), tt.lang).Return(tt.findAllResp, tt.findAllErr).AnyTimes()

			u := NewMuscleUsecase(mockMuscleRepository)

			_, err := u.ListMuscles(tt.ctx, tt.lang)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
