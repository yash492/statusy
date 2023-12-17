package types

type NullType[T any] struct {
	IsValid bool
	Value   T
}
