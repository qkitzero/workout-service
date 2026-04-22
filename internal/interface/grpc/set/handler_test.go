package set

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	setv1 "github.com/qkitzero/workout-service/gen/go/set/v1"
	"github.com/qkitzero/workout-service/internal/domain/set"
	mocksappset "github.com/qkitzero/workout-service/mocks/application/set"
	mocksset "github.com/qkitzero/workout-service/mocks/domain/set"
)

func TestCreateSet(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		success      bool
		ctx          context.Context
		exercise     string
		rep          int32
		weight       float64
		trainedAt    time.Time
		createSetErr error
	}{
		{"success create set", true, context.Background(), "bench press", 10, 60.0, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), nil},
		{"failure create set error", false, context.Background(), "bench press", 10, 60.0, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), fmt.Errorf("create set error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSetUsecase := mocksappset.NewMockSetUsecase(ctrl)
			mockSet := mocksset.NewMockSet(ctrl)
			mockSetUsecase.EXPECT().CreateSet(tt.ctx, tt.exercise, tt.rep, tt.weight, tt.trainedAt).Return(mockSet, tt.createSetErr).AnyTimes()
			mockSetID := set.NewSetID()
			mockSet.EXPECT().ID().Return(mockSetID).AnyTimes()

			setHandler := NewSetHandler(mockSetUsecase)

			req := &setv1.CreateSetRequest{
				Exercise:  tt.exercise,
				Rep:       tt.rep,
				Weight:    tt.weight,
				TrainedAt: timestamppb.New(tt.trainedAt),
			}

			_, err := setHandler.CreateSet(tt.ctx, req)
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
		name        string
		success     bool
		ctx         context.Context
		listSetsErr error
	}{
		{"success list sets", true, context.Background(), nil},
		{"failure list sets error", false, context.Background(), fmt.Errorf("list sets error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSetUsecase := mocksappset.NewMockSetUsecase(ctrl)
			mockSet := mocksset.NewMockSet(ctrl)
			mockSetID := set.NewSetID()
			mockSet.EXPECT().ID().Return(mockSetID).AnyTimes()
			mockSet.EXPECT().Exercise().Return(set.Exercise("bench press")).AnyTimes()
			mockSet.EXPECT().Rep().Return(set.Rep(10)).AnyTimes()
			mockSet.EXPECT().Weight().Return(set.Weight(60.0)).AnyTimes()
			mockSet.EXPECT().TrainedAt().Return(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)).AnyTimes()
			mockSet.EXPECT().CreatedAt().Return(time.Now()).AnyTimes()
			mockSetUsecase.EXPECT().ListSets(tt.ctx).Return([]set.Set{mockSet}, tt.listSetsErr).AnyTimes()

			setHandler := NewSetHandler(mockSetUsecase)

			_, err := setHandler.ListSets(tt.ctx, &setv1.ListSetsRequest{})
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
