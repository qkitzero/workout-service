package user

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"

	userv1 "github.com/qkitzero/user-service/gen/go/user/v1"
	mocks "github.com/qkitzero/workout-service/mocks/external/user/v1"
)

func TestGetUser(t *testing.T) {
	t.Parallel()
	accessToken := "accessToken"
	tests := []struct {
		name       string
		success    bool
		ctx        context.Context
		getUserErr error
	}{
		{
			name:       "success get user",
			success:    true,
			ctx:        metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+accessToken)),
			getUserErr: nil,
		},
		{
			name:       "failure missing metadata",
			success:    false,
			ctx:        context.Background(),
			getUserErr: nil,
		},
		{
			name:       "failure get user error",
			success:    false,
			ctx:        metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+accessToken)),
			getUserErr: errors.New("get user error"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockUserServiceClient(ctrl)
			mockGetUserResponse := &userv1.GetUserResponse{
				UserId: "fe8c2263-bbac-4bb9-a41d-b04f5afc4425",
			}
			mockClient.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(mockGetUserResponse, tt.getUserErr).AnyTimes()

			userService := NewUserService(mockClient)

			_, err := userService.GetUser(tt.ctx)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
