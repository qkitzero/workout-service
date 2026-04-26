package muscle

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	musclev1 "github.com/qkitzero/workout-service/gen/go/muscle/v1"
	"github.com/qkitzero/workout-service/internal/domain/i18n"
	"github.com/qkitzero/workout-service/internal/domain/muscle"
	mocksappmuscle "github.com/qkitzero/workout-service/mocks/application/muscle"
)

func TestListMuscles(t *testing.T) {
	t.Parallel()

	code, err := muscle.NewCode("chest")
	if err != nil {
		t.Fatalf("failed to new code: %v", err)
	}
	name, err := muscle.NewName("胸")
	if err != nil {
		t.Fatalf("failed to new name: %v", err)
	}
	sample := muscle.NewMuscle(muscle.NewMuscleID(), code, name)

	tests := []struct {
		name           string
		ctx            context.Context
		lang           string
		callUsecase    bool
		wantLang       i18n.Language
		listMusclesErr error
		wantCode       codes.Code
	}{
		{"success list muscles", context.Background(), "ja", true, i18n.LanguageJa, nil, codes.OK},
		{"success default lang", context.Background(), "", true, i18n.LanguageJa, nil, codes.OK},
		{"failure invalid lang", context.Background(), "JA", false, "", nil, codes.InvalidArgument},
		{"failure list muscles error", context.Background(), "ja", true, i18n.LanguageJa, fmt.Errorf("list muscles error"), codes.Internal},
		{"failure status preserved", context.Background(), "ja", true, i18n.LanguageJa, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappmuscle.NewMockMuscleUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().ListMuscles(tt.ctx, tt.wantLang).Return([]muscle.Muscle{sample}, tt.listMusclesErr).Times(1)
			}

			handler := NewMuscleHandler(mockUsecase)

			req := &musclev1.ListMusclesRequest{Lang: tt.lang}
			_, err := handler.ListMuscles(tt.ctx, req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}
