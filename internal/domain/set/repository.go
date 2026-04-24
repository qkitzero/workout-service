package set

import (
	"context"

	"github.com/qkitzero/workout-service/internal/domain/user"
)

type SetRepository interface {
	Create(ctx context.Context, set Set) error
	FindByUserID(ctx context.Context, userID user.UserID) ([]Set, error)
}
