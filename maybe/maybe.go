package maybe

type (
	Maybe[T any] struct {
		value *T
	}
)

func Just[T any](value T) Maybe[T] {
	return Maybe[T]{
		value: &value,
	}
}

func Nothing[T any]() Maybe[T] {
	return Maybe[T]{
		value: nil,
	}
}

// A nil pointer indicates `Nothing`, conversely a non-nil pointer indicates a `Just` value.
func FromGoPointer[T any](ptr *T) Maybe[T] {
	return Maybe[T]{
		value: ptr,
	}
}

func Bind[T1 any, T2 any](m Maybe[T1], f func(T1) Maybe[T2]) Maybe[T2] {
	if m.IsNothing() {
		return Nothing[T2]()
	}
	return f(m.Value())
}

func (m Maybe[T]) IsJust() bool {
	return m.value != nil
}

func (m Maybe[T]) IsNothing() bool {
	return m.value == nil
}

func (m Maybe[T]) ToGoPointer() *T {
	return m.value
}

func (m Maybe[T]) Value() T {
	return *m.value
}

func (m Maybe[T]) OrValue(defaultValue T) T {
	if m.IsNothing() {
		return defaultValue
	}
	return m.Value()
}
