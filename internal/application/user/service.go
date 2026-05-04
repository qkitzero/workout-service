package user

import "context"

type UserService interface {
	GetUser(ctx context.Context) (string, error)
}
