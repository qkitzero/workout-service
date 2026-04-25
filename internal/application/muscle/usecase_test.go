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
		lang        string
		findAllResp []muscle.Muscle
		findAllErr  error
		wantLang    i18n.Language
	}{
		{"success default lang", true, context.Background(), "", []muscle.Muscle{sample}, nil, i18n.LanguageJa},
		{"success ja", true, context.Background(), "ja", []muscle.Muscle{sample}, nil, i18n.LanguageJa},
		{"success en", true, context.Background(), "en", []muscle.Muscle{sample}, nil, i18n.Language("en")},
		{"success empty result", true, context.Background(), "ja", []muscle.Muscle{}, nil, i18n.LanguageJa},
		{"failure invalid lang", false, context.Background(), "JA", nil, nil, ""},
		{"failure find all error", false, context.Background(), "ja", nil, errors.New("find all error"), i18n.LanguageJa},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockMuscleRepository := mocksmuscle.NewMockMuscleRepository(ctrl)
			if tt.wantLang != "" {
				mockMuscleRepository.EXPECT().FindAll(gomock.Any(), tt.wantLang).Return(tt.findAllResp, tt.findAllErr).AnyTimes()
			}

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
