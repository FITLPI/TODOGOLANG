package domain

// Nullable — трёхсостоящее значение:
//   - поле не передано (Set == false) -> оставить как было;
//   - передан явный null (Set == true, Value == nil) -> очистить;
//   - передано значение (Set == true, Value != nil) -> записать.
type Nullable[T any] struct {
	Set   bool
	Value *T
}

func NewNullableNotSet[T any]() Nullable[T] {
	return Nullable[T]{Set: false, Value: nil}
}

func NewNullableNull[T any]() Nullable[T] {
	return Nullable[T]{Set: true, Value: nil}
}

func NewNullableValue[T any](v T) Nullable[T] {
	value := v
	return Nullable[T]{Set: true, Value: &value}
}
