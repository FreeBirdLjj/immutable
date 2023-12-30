package list

import (
	"github.com/freebirdljj/immutable/comparator"
	immutable_func "github.com/freebirdljj/immutable/func"
)

type (
	List[T any] struct {
		value T
		next  *List[T]
	}
)

func Cons[T any](x T, xs *List[T]) *List[T] {
	return &List[T]{
		value: x,
		next:  xs,
	}
}

func FromGoSlice[T any](xs []T) *List[T] {
	l := (*List[T])(nil)
	for i := range xs {
		l = Cons(xs[len(xs)-1-i], l)
	}
	return l
}

func Repeat[T any](x T) *List[T] {
	xs := List[T]{
		value: x,
	}
	xs.next = &xs
	return &xs
}

// CAUTION: `xs` can't be nil.
func Cycle[T any](xs *List[T]) *List[T] {
	first := Cons(xs.value, nil)
	return maplist(xs, func(p *List[T]) *List[T] {
		if p == nil || p == xs {
			return first
		}
		return Cons(p.value, nil)
	})
}

func Map[T1 any, T2 any](xs *List[T1], f func(T1) T2) *List[T2] {
	return maplist(xs, func(p *List[T1]) *List[T2] {
		if p == nil {
			return nil
		}
		return Cons(f(p.value), nil)
	})
}

// CAUTION: Only invoke `Foldl` with finite list `xs`.
func Foldl[T1 any, T2 any](xs *List[T1], init T2, f func(acc T2, x T1) T2) T2 {
	res := init
	for p := xs; p != nil; p = p.next {
		res = f(res, p.value)
	}
	return res
}

// CAUTION: Only invoke `Foldr` with finite list `xs`.
func Foldr[T1 any, T2 any](xs *List[T1], init T2, f func(x T1, acc T2) T2) T2 {
	return Foldl(
		xs,
		immutable_func.Identity[T2],
		func(folded func(acc T2) T2, value T1) func(acc T2) T2 {
			return func(acc T2) T2 {
				return folded(f(value, acc))
			}
		},
	)(init)
}

func Concat[T any](xss *List[*List[T]]) *List[T] {
	return maplist(xss, func(p *List[*List[T]]) *List[T] {
		if p == nil {
			return nil
		}
		return p.value.clone()
	})
}

// NOTE: The `next` field of the last node of the list returned by `f` may be modified
func maplist[T1 any, T2 any](xs *List[T1], f func(*List[T1]) *List[T2]) *List[T2] {

	head := List[T2]{}
	prev := &head
	nodeMap := make(map[*List[T1]]*List[T2])

	for p := xs; p != nil; p = p.next {

		if mappedNode, mapped := nodeMap[p]; mapped {
			// NOTE: For the second lap run, skip all nodes mapped to `nil`.
			if mappedNode != nil {
				prev.next = mappedNode
				return head.next
			}
			circleEntry := p
			for p = p.next; p != circleEntry; p = p.next {
				if mappedNode := nodeMap[p]; mappedNode != nil {
					prev.next = nodeMap[p]
					return head.next
				}
			}
			return head.next
		}

		newNode := f(p)
		nodeMap[p] = newNode
		prev.next = newNode

		if !prev.isFinite() {
			return head.next
		}

		for prev.next != nil {
			prev = prev.next
		}
	}

	prev.next = f(nil)
	return head.next
}

// CAUTION: `xs` can't be nil.
func (xs *List[T]) Uncons() (value T, next *List[T]) {
	return xs.value, xs.next
}

func (xs *List[T]) Empty() bool {
	return xs == nil
}

// CAUTION: Only invoke `Length()` with finite list.
func (xs *List[T]) Length() int {

	if xs == nil {
		return 0
	}

	n := 0
	for p := xs; p != nil; p = p.next {
		n++
	}
	return n
}

func (xs *List[T]) Append(ys *List[T]) *List[T] {

	if !xs.isFinite() || ys == nil {
		return xs
	}

	return maplist(xs, func(p *List[T]) *List[T] {
		if p == nil {
			return ys
		}
		return Cons(p.value, nil)
	})
}

func (xs *List[T]) Take(n int) *List[T] {
	res := make([]T, 0, n)
	for p := xs; p != nil && n > 0; p = p.next {
		res = append(res, p.value)
		n--
	}
	return FromGoSlice(res)
}

func (xs *List[T]) Drop(n int) *List[T] {
	p := xs
	for p != nil && n > 0 {
		n--
		p = p.next
	}
	return p
}

func (xs *List[T]) Find(predicate func(T) bool) *T {
	visited := make(map[*List[T]]bool)
	for p := xs; p != nil; p = p.next {

		if visited[p] {
			return nil
		}
		visited[p] = true

		if predicate(p.value) {
			return &p.value
		}
	}
	return nil
}

func (xs *List[T]) Filter(predicate func(T) bool) *List[T] {
	return maplist(xs, func(p *List[T]) *List[T] {
		if p == nil || !predicate(p.value) {
			return nil
		}
		return Cons(p.value, nil)
	})
}

func (xs *List[T]) Sort(cmp comparator.Comparator[T]) *List[T] {

	if xs == nil || xs.next == nil || xs.next == xs {
		return xs
	}

	head, tail := xs.Uncons()
	lessPart := tail.Filter(func(x T) bool { return cmp(x, head) < 0 }).Sort(cmp)

	if !lessPart.isFinite() {
		return lessPart
	}

	return lessPart.Append(Cons(head, tail.Filter(func(x T) bool { return cmp(x, head) >= 0 }).Sort(cmp)))
}

func (xs *List[T]) Reverse() *List[T] {
	res := (*List[T])(nil)
	for p := xs; p != nil; p = p.next {
		res = Cons(p.value, res)
	}
	return res
}

func (xs *List[T]) Intersperse(sep T) *List[T] {
	return maplist(xs, func(p *List[T]) *List[T] {
		if p == nil || p.next == nil {
			return p
		}
		return &List[T]{
			value: p.value,
			next: &List[T]{
				value: sep,
			},
		}
	})
}

func (xs *List[T]) IsIsomorphicTo(ys *List[T], cmp comparator.Comparator[T]) bool {

	xVisited := map[*List[T]]bool{nil: true}
	yVisited := map[*List[T]]bool{nil: true}

	for !xVisited[xs] || !yVisited[ys] {

		if xs == nil || ys == nil || cmp(xs.value, ys.value) != 0 {
			return false
		}

		xVisited[xs] = true
		yVisited[ys] = true

		xs = xs.next
		ys = ys.next
	}

	return true
}

func (xs *List[T]) All(predicate func(T) bool) bool {
	visited := make(map[*List[T]]bool)
	for p := xs; p != nil; p = p.next {

		if visited[p] {
			return true
		}
		visited[p] = true

		if !predicate(p.value) {
			return false
		}
	}
	return true
}

func (xs *List[T]) Any(predicate func(T) bool) bool {
	visited := make(map[*List[T]]bool)
	for p := xs; p != nil; p = p.next {

		if visited[p] {
			return false
		}
		visited[p] = true

		if predicate(p.value) {
			return true
		}
	}
	return false
}

// CAUTION: Only invoke `ToGoSlice()` with finite list.
func (xs *List[T]) ToGoSlice() []T {

	if xs == nil {
		return nil
	}

	res := make([]T, 0, xs.Length())
	for p := xs; p != nil; p = p.next {
		res = append(res, p.value)
	}
	return res
}

func (xs *List[T]) clone() *List[T] {
	return maplist(xs, func(p *List[T]) *List[T] {
		if p == nil {
			return nil
		}
		return Cons(p.value, nil)
	})
}

func (xs *List[T]) isFinite() bool {

	if xs == nil {
		return true
	}

	for pSlow, pFast := xs, xs.next; pSlow != pFast; pSlow = pSlow.next {
		for times := 0; times < 2; times++ {
			if pFast == nil {
				return true
			}
			pFast = pFast.next
		}
	}

	return false
}
