package exercise

import "fmt"

type Category string

const (
	CategoryCompound  Category = "compound"
	CategoryIsolation Category = "isolation"
)

func (c Category) String() string {
	return string(c)
}

func NewCategory(s string) (Category, error) {
	switch Category(s) {
	case CategoryCompound, CategoryIsolation:
		return Category(s), nil
	default:
		return Category(""), fmt.Errorf("invalid category: %q", s)
	}
}
