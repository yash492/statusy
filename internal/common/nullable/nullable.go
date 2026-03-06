package nullable

type Nullable[T any] struct {
	value T
	valid bool
}

func (n Nullable[T]) Get() (T, bool) {
	return n.value, n.valid
}

func SetValue[T any](value T, isValid bool) Nullable[T] {
	return Nullable[T]{
		value: value,
		valid: isValid,
	}
}
