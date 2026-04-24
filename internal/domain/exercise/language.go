package exercise

import (
	"fmt"
	"regexp"
)

type Language string

const LanguageJa Language = "ja"

var languagePattern = regexp.MustCompile(`^[a-z]{2}(-[A-Z]{2})?$`)

func (l Language) String() string {
	return string(l)
}

func NewLanguage(s string) (Language, error) {
	if !languagePattern.MatchString(s) {
		return Language(""), fmt.Errorf("invalid language: %q", s)
	}
	return Language(s), nil
}
