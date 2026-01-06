package common

type Nullable[T any] struct {
	Value T
	Valid bool
}

func (n Nullable[T]) Get() (T, bool) {
	return n.Value, n.Valid
}

func SetNullableValue[T any](value T) Nullable[T] {
	return Nullable[T]{
		Value: value,
		Valid: true,
	}
}
