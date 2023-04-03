package list

import (
	"github.com/freebirdljj/immutable/comparator"
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

func FromSlice[T any](xs []T) *List[T] {
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

	first := *xs
	nodeMap := map[*List[T]]*List[T]{
		xs:  &first,
		nil: &first,
	}

	for p := &first; ; p = p.next {
		if mappedNode, mapped := nodeMap[p.next]; mapped {
			p.next = mappedNode
			return &first
		}
		newNode := *p.next
		nodeMap[p.next] = &newNode
		p.next = &newNode
	}
}

func Map[T1 any, T2 any](xs *List[T1], f func(T1) T2) *List[T2] {

	head := List[T2]{}
	prev := &head
	nodeMap := map[*List[T1]]*List[T2]{
		nil: nil,
	}

	for p := xs; ; p = p.next {
		if mappedNode, mapped := nodeMap[p]; mapped {
			prev.next = mappedNode
			return head.next
		}
		newNode := List[T2]{
			value: f(p.value),
		}
		nodeMap[p] = &newNode
		prev.next = &newNode
		prev = &newNode
	}
}

// CAUTION: Only invoke `Foldl` with finite list `xs`.
func Foldl[T1 any, T2 any](xs *List[T1], init T2, f func(acc T2, x T1) T2) T2 {
	res := init
	for p := xs; p != nil; p = p.next {
		res = f(res, p.value)
	}
	return res
}

// CAUTION: `xs` can't be nil.
func (xs *List[T]) Uncons() (value T, next *List[T]) {
	return xs.value, xs.next
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

	head := List[T]{
		next: xs,
	}
	nodeMap := map[*List[T]]*List[T]{
		nil: ys,
	}

	for p := &head; ; p = p.next {
		if mappedNode, mapped := nodeMap[p.next]; mapped {
			p.next = mappedNode
			return head.next
		}
		newNode := *p.next
		nodeMap[p.next] = &newNode
		p.next = &newNode
	}
}

func (xs *List[T]) Take(n int) *List[T] {
	res := make([]T, 0, n)
	for p := xs; p != nil && n > 0; p = p.next {
		res = append(res, p.value)
		n--
	}
	return FromSlice(res)
}

func (xs *List[T]) Drop(n int) *List[T] {
	p := xs
	for p != nil && n > 0 {
		n--
		p = p.next
	}
	return p
}

func (xs *List[T]) Filter(predicate func(T) bool) *List[T] {

	nodeMap := map[*List[T]]*List[T]{
		nil: nil,
	}

	head := List[T]{}
	prev := &head

	for p := xs; p != nil; p = p.next {

		if mappedNode, mapped := nodeMap[p]; mapped {
			// NOTE: For the second lap run, avoid repeatedly running the predicate at the same node.
			//       predicate(p.value) <=> nodeMap[p] != nil
			if mappedNode != nil {
				prev.next = mappedNode
				break
			}
			circleEntry := p
			for p = p.next; p != circleEntry; p = p.next {
				if mappedNode := nodeMap[p]; mappedNode != nil {
					prev.next = nodeMap[p]
					break
				}
			}
			break
		}

		if predicate(p.value) {
			pCopy := List[T]{
				value: p.value,
			}
			nodeMap[p] = &pCopy
			prev.next = &pCopy
			prev = &pCopy
		} else {
			nodeMap[p] = nil
		}
	}

	return head.next
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

// CAUTION: Only invoke `ToSlice()` with finite list.
func (xs *List[T]) ToSlice() []T {

	if xs == nil {
		return nil
	}

	res := make([]T, 0, xs.Length())
	for p := xs; p != nil; p = p.next {
		res = append(res, p.value)
	}
	return res
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
