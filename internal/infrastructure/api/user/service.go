package user

import (
	"context"
	"errors"

	userv1 "github.com/qkitzero/user-service/gen/go/user/v1"
	"github.com/qkitzero/workout-service/internal/application/user"
	"google.golang.org/grpc/metadata"
)

type userService struct {
	client userv1.UserServiceClient
}

func NewUserService(client userv1.UserServiceClient) user.UserService {
	return &userService{client: client}
}

func (s *userService) GetUser(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata is missing")
	}

	ctx = metadata.NewOutgoingContext(ctx, md)

	getUserRequest := &userv1.GetUserRequest{}

	getUserResponse, err := s.client.GetUser(ctx, getUserRequest)
	if err != nil {
		return "", err
	}

	return getUserResponse.GetUserId(), nil
}
