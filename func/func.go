package immutable_func

func Identity[T any](x T) T {
	return x
}

func Zero[T any]() T {
	var x T
	return x
}
