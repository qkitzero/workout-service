package set

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/qkitzero/workout-service/internal/domain/set"
	mocksappauth "github.com/qkitzero/workout-service/mocks/application/auth"
	mocks "github.com/qkitzero/workout-service/mocks/domain/set"
)

func TestCreateSet(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		success        bool
		ctx            context.Context
		userID         string
		exercise       string
		rep            int32
		weight         float64
		trainedAt      time.Time
		verifyTokenErr error
		createErr      error
	}{
		{"success create set", true, context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "bench press", 10, 60.0, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), nil, nil},
		{"failure verify token error", false, context.Background(), "", "bench press", 10, 60.0, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), fmt.Errorf("verify token error"), nil},
		{"failure empty user id", false, context.Background(), "", "bench press", 10, 60.0, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), nil, nil},
		{"failure empty exercise", false, context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "", 10, 60.0, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), nil, nil},
		{"failure invalid rep", false, context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "bench press", 0, 60.0, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), nil, nil},
		{"failure negative weight", false, context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "bench press", 10, -1.0, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), nil, nil},
		{"failure create error", false, context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "bench press", 10, 60.0, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), nil, errors.New("create error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockSetRepository := mocks.NewMockSetRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockSetRepository.EXPECT().Create(gomock.Any()).Return(tt.createErr).AnyTimes()

			setUsecase := NewSetUsecase(mockAuthService, mockSetRepository)

			_, err := setUsecase.CreateSet(tt.ctx, tt.exercise, tt.rep, tt.weight, tt.trainedAt)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestListSets(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		success         bool
		ctx             context.Context
		userID          string
		verifyTokenErr  error
		findByUserIDErr error
	}{
		{"success list sets", true, context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", nil, nil},
		{"failure verify token error", false, context.Background(), "", fmt.Errorf("verify token error"), nil},
		{"failure empty user id", false, context.Background(), "", nil, nil},
		{"failure find by user id error", false, context.Background(), "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", nil, errors.New("find by user id error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockSetRepository := mocks.NewMockSetRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockSetRepository.EXPECT().FindByUserID(gomock.Any()).Return([]set.Set{}, tt.findByUserIDErr).AnyTimes()

			setUsecase := NewSetUsecase(mockAuthService, mockSetRepository)

			_, err := setUsecase.ListSets(tt.ctx)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
