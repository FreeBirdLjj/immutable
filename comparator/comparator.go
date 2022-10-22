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
