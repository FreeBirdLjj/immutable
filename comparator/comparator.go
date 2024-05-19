package comparator

import (
	"golang.org/x/exp/constraints"
)

type (
	// A negative return value indicates l < r;
	// a return value of 0 indicates l == r;
	// a positive return value indicates l > r.
	Comparator[T any] func(l T, r T) int
)

func OrderedComparator[T constraints.Ordered](l T, r T) int {
	switch {
	case l < r:
		return -1
	case l > r:
		return 1
	default: // <- case l == r, but just make compiler happy
		return 0
	}
}

func CascadeComparator[T1 any, T2 any](cmp Comparator[T1], mapping func(T2) T1) Comparator[T2] {
	return func(l T2, r T2) int {
		return cmp(mapping(l), mapping(r))
	}
}
