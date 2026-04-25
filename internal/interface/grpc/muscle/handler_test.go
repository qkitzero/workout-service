package muscle

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"

	musclev1 "github.com/qkitzero/workout-service/gen/go/muscle/v1"
	"github.com/qkitzero/workout-service/internal/domain/muscle"
	mocksappmuscle "github.com/qkitzero/workout-service/mocks/application/muscle"
)

func TestListMuscles(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		success        bool
		ctx            context.Context
		lang           string
		listMusclesErr error
	}{
		{"success list muscles", true, context.Background(), "ja", nil},
		{"success default lang", true, context.Background(), "", nil},
		{"failure list muscles error", false, context.Background(), "ja", fmt.Errorf("list muscles error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappmuscle.NewMockMuscleUsecase(ctrl)

			code, err := muscle.NewCode("chest")
			if err != nil {
				t.Errorf("failed to new code: %v", err)
			}
			name, err := muscle.NewName("胸")
			if err != nil {
				t.Errorf("failed to new name: %v", err)
			}
			sample := muscle.NewMuscle(muscle.NewMuscleID(), code, name)

			mockUsecase.EXPECT().ListMuscles(tt.ctx, tt.lang).Return([]muscle.Muscle{sample}, tt.listMusclesErr).AnyTimes()

			handler := NewMuscleHandler(mockUsecase)

			req := &musclev1.ListMusclesRequest{Lang: tt.lang}
			_, err = handler.ListMuscles(tt.ctx, req)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
