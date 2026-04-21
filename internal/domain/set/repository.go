package set

type SetRepository interface {
	Create(set Set) error
}
