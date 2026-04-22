package set

import "github.com/qkitzero/workout-service/internal/domain/user"

type SetRepository interface {
	Create(set Set) error
	FindByUserID(userID user.UserID) ([]Set, error)
}
