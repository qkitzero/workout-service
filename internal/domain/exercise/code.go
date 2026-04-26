package exercise

import (
	"fmt"
	"regexp"
)

type Code string

var codePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

func (c Code) String() string {
	return string(c)
}

func NewCode(s string) (Code, error) {
	if !codePattern.MatchString(s) {
		return Code(""), fmt.Errorf("invalid code: %q", s)
	}
	return Code(s), nil
}
