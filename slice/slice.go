package slice

import (
	"sort"

	"github.com/freebirdljj/immutable/comparator"
	"github.com/freebirdljj/immutable/maybe"
)

type (
	Slice[T any] []T
)

func FromGoSlice[T any](xs []T) Slice[T] {
	return Slice[T](xs)
}

func Map[T1 any, T2 any](xs Slice[T1], f func(T1) T2) Slice[T2] {
	ys := make(Slice[T2], len(xs))
	for i, x := range xs {
		ys[i] = f(x)
	}
	return ys
}

func Foldl[T1 any, T2 any](xs Slice[T1], init T2, f func(acc T2, x T1) T2) T2 {
	res := init
	for _, x := range xs {
		res = f(res, x)
	}
	return res
}

func Foldr[T1 any, T2 any](xs Slice[T1], init T2, f func(x T1, acc T2) T2) T2 {
	res := init
	for i := range xs {
		x := xs[len(xs)-1-i]
		res = f(x, res)
	}
	return res
}

// CAUTION: Only invoke `MaximumBy` with non-empty slice `xs`.
func MaximumBy[T any](xs Slice[T], cmp comparator.Comparator[T]) T {
	max := xs[0]
	for _, x := range xs[1:] {
		if cmp(max, x) < 0 {
			max = x
		}
	}
	return max
}

// CAUTION: Only invoke `MinimumBy` with non-empty slice `xs`.
func MinimumBy[T any](xs Slice[T], cmp comparator.Comparator[T]) T {
	min := xs[0]
	for _, x := range xs[1:] {
		if cmp(min, x) > 0 {
			min = x
		}
	}
	return min
}

func GroupBy[T any](xs Slice[T], cmp comparator.Comparator[T]) Slice[Slice[T]] {
	res := Slice[Slice[T]](nil)
	for len(xs) > 0 {
		head := xs[0]
		headEqs, rest := xs.Partition(func(x T) bool { return cmp(head, x) == 0 })
		res = append(res, headEqs)
		xs = rest
	}
	return res
}

func Concat[T any](xss Slice[Slice[T]]) Slice[T] {
	switch len(xss) {
	case 0:
		return nil
	case 1:
		return xss[0]
	default:
		cnt := Foldl(xss, 0, func(acc int, xs Slice[T]) int { return acc + len(xs) })
		res := make(Slice[T], 0, cnt)
		for _, xs := range xss {
			res = append(res, xs...)
		}
		return res
	}
}

func (xs Slice[T]) Empty() bool {
	return len(xs) == 0
}

func (xs Slice[T]) Tail() Slice[T] {
	if len(xs) == 0 {
		return nil
	}
	return xs[1:]
}

func (xs Slice[T]) Append(elems ...T) Slice[T] {

	if len(elems) == 0 {
		return xs
	}

	if len(xs) == 0 {
		return elems
	}

	res := make(Slice[T], len(xs)+len(elems))
	copy(res, xs)
	copy(res[len(xs):], elems)
	return res
}

func (xs Slice[T]) Take(n int) Slice[T] {
	return xs[:n]
}

func (xs Slice[T]) Drop(n int) Slice[T] {
	return xs[n:]
}

func (xs Slice[T]) Find(predicate func(T) bool) maybe.Maybe[T] {
	for i, x := range xs {
		if predicate(x) {
			return maybe.FromGoPointer(&xs[i])
		}
	}
	return maybe.Nothing[T]()
}

func (xs Slice[T]) Filter(predicate func(T) bool) Slice[T] {
	res := make(Slice[T], 0, len(xs))
	for _, x := range xs {
		if predicate(x) {
			res = append(res, x)
		}
	}
	return res
}

func (xs Slice[T]) Partition(predicate func(T) bool) (satisfied Slice[T], unsatisfied Slice[T]) {

	satisfied = make(Slice[T], 0, len(xs))
	unsatisfied = make(Slice[T], 0, len(xs))

	for _, x := range xs {
		if predicate(x) {
			satisfied = append(satisfied, x)
		} else {
			unsatisfied = append(unsatisfied, x)
		}
	}

	return satisfied, unsatisfied
}

func (xs Slice[T]) Sort(cmp comparator.Comparator[T]) Slice[T] {

	lessMaker := func(s Slice[T]) func(i int, j int) bool {
		return func(i int, j int) bool {
			return cmp(s[i], s[j]) < 0
		}
	}

	if len(xs) <= 1 || sort.SliceIsSorted(xs, lessMaker(xs)) {
		return xs
	}

	res := make(Slice[T], len(xs))
	copy(res, xs)
	sort.Slice(res, lessMaker(res))
	return res
}

func (xs Slice[T]) Reverse() Slice[T] {
	res := make(Slice[T], len(xs))
	for i := range xs {
		x := xs[len(xs)-1-i]
		res[i] = x
	}
	return res
}

func (xs Slice[T]) Intersperse(sep T) Slice[T] {
	if len(xs) <= 1 {
		return xs
	}
	res := make(Slice[T], len(xs)*2-1)
	res[0] = xs[0]
	for i, x := range xs.Tail() {
		res[i*2+1] = sep
		res[i*2+2] = x
	}
	return res
}

func (xs Slice[T]) All(predicate func(T) bool) bool {
	for _, x := range xs {
		if !predicate(x) {
			return false
		}
	}
	return true
}

func (xs Slice[T]) Any(predicate func(T) bool) bool {
	for _, x := range xs {
		if predicate(x) {
			return true
		}
	}
	return false
}

func (xs Slice[T]) ToGoSlice() []T {
	return []T(xs)
}
