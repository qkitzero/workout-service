package user

import "fmt"

type UserID string

func (u UserID) String() string {
	return string(u)
}

func NewUserID(s string) (UserID, error) {
	if s == "" {
		return UserID(""), fmt.Errorf("user id is empty")
	}
	return UserID(s), nil
}
