package workout

import "errors"

var (
	ErrWorkoutNotFound        = errors.New("workout not found")
	ErrWorkoutAlreadyFinished = errors.New("workout already finished")
	ErrWorkoutForbidden       = errors.New("workout does not belong to user")
)
