package common

type Nullable[T any] struct {
	Value T
	Valid bool
}

func (n Nullable[T]) Get() (T, bool) {
	return n.Value, n.Valid
}
