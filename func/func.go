package immutable_func

func Identity[T any](x T) T {
	return x
}

func Konst[T1 any, T2 any](k T2) func(T1) T2 {
	return func(T1) T2 {
		return k
	}
}

func Zero[T any]() T {
	var x T
	return x
}
