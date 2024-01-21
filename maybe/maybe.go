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
