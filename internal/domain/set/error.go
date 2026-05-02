package set

import "errors"

var (
	ErrSetNotFound  = errors.New("set not found")
	ErrSetForbidden = errors.New("set does not belong to user")
)
