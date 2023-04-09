package list

import (
	"math"
	"reflect"
	"sort"
	"strconv"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/require"

	"github.com/freebirdljj/immutable/comparator"
	immutable_func "github.com/freebirdljj/immutable/func"
)

func TestCycle(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"cycle(xs) is infinite": func(xs []int, last int) bool {
			nonemptySlice := append(xs, last)
			xl := FromSlice(nonemptySlice)
			return !Cycle(xl).isFinite()
		},
		"cycle(xs).take(2 * xs.length()) == xs ++ xs": func(xs []int, last int) bool {
			nonemptySlice := append(xs, last)
			xl := FromSlice(nonemptySlice)
			return slicesEqual(Cycle(xl).Take(2*len(nonemptySlice)).ToSlice(), append(nonemptySlice, nonemptySlice...))
		},
		"cycle of an infinite list will be the same": func(prefixes []int, xs []int, last int) bool {
			nonemptySlice := append(xs, last)
			xl := Cycle(FromSlice(nonemptySlice))
			sampleLen := len(prefixes) + 2*len(nonemptySlice)
			return slicesEqual(Cycle(xl).Take(sampleLen).ToSlice(), xl.Take(sampleLen).ToSlice())
		},
	})
}

func TestMap(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.map(f).length() == xs.length()": func(xs []int) bool {
			f := func(x int) int { return x + 1 }
			l := FromSlice(xs)
			return Map(l, f).Length() == l.Length()
		},
		"cycle(xs).map(f) == cycle(xs.map(f))": func(xs []int, last int) bool {
			f := func(x int) int { return x + 1 }
			nonemptySlice := append(xs, last)
			xl := FromSlice(nonemptySlice)
			sampleLen := 2 * (len(xs) + 1)
			return slicesEqual(Map(Cycle(xl), f).Take(sampleLen).ToSlice(), Cycle(Map(xl, f)).Take(sampleLen).ToSlice())
		},
		"xs.map(f1).map(f2) == xs.map(f2 . f1)": func(xs []int) bool {
			f1 := func(x int) int { return x + 1 }
			f2 := strconv.Itoa
			xl := FromSlice(xs)
			return slicesEqual(Map(Map(xl, f1), f2).ToSlice(), Map(xl, func(x int) string { return f2(f1(x)) }).ToSlice())
		},
	})
}

func TestFoldl(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.foldl([], flip(cons)).reverse() == xs": func(xs []int) bool {
			xl := FromSlice(xs)
			return slicesEqual(Foldl(xl, nil, func(acc *List[int], x int) *List[int] { return Cons(x, acc) }).Reverse().ToSlice(), xs)
		},
	})
}

func TestFoldr(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.foldr([], cons) == xs": func(xs []int) bool {
			xl := FromSlice(xs)
			return slicesEqual(Foldr(xl, nil, Cons[int]).ToSlice(), xs)
		},
	})
}

func TestListAppend(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.append(ys) == xs ++ ys": func(xs []int, ys []int) bool {
			xl := FromSlice(xs)
			yl := FromSlice(ys)
			return slicesEqual(xl.Append(yl).ToSlice(), append(xs, ys...))
		},
		"cycle(xs).append(ys) == cycle(xs)": func(xs []int, ys []int, xLast int) bool {
			nonemptySliceX := append(xs, xLast)
			xl := FromSlice(nonemptySliceX)
			yl := FromSlice(ys)
			return Cycle(xl).Append(yl).IsIsomorphicTo(Cycle(xl), comparator.OrderedComparator[int])
		},
	})
}

func TestListTake(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.append(ys).take(xs.length()) == xs": func(xs []int, ys []int) bool {
			xl := FromSlice(xs)
			yl := FromSlice(ys)
			return slicesEqual(xl.Append(yl).Take(len(xs)).ToSlice(), xs)
		},
		"xs.take(n) == xs if n >= xs.length()": func(xs []int, delta uint8) bool {
			n := len(xs) + int(delta)
			xl := FromSlice(xs)
			return slicesEqual(xl.Take(n).ToSlice(), xs)
		},
	})
}

func TestListDrop(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.append(ys).drop(xs.length()) == ys": func(xs []int, ys []int) bool {
			xl := FromSlice(xs)
			yl := FromSlice(ys)
			return slicesEqual(xl.Append(yl).Drop(len(xs)).ToSlice(), ys)
		},
		"xs.drop(n) == nil if n >= xs.length()": func(xs []int, delta uint8) bool {
			n := len(xs) + int(delta)
			xl := FromSlice(xs)
			return xl.Drop(n) == nil
		},
	})
}

func TestListFilter(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.filter(p).append(ys.filter(p)) == xs.append(ys).filter(p)": func(xs []int, ys []int) bool {
			predicate := func(x int) bool { return x%2 == 0 }
			xl := FromSlice(xs)
			yl := FromSlice(ys)
			return reflect.DeepEqual(xl.Filter(predicate).Append(yl.Filter(predicate)).ToSlice(), xl.Append(yl).Filter(predicate).ToSlice())
		},
		"xs.filter(konst(false)) == nil": func(xs []int) bool {
			predicate := immutable_func.Konst[int](false)
			xl := FromSlice(xs)
			return xl.Filter(predicate) == nil
		},
		"cycle(xs).filter(konst(false)) == nil": func(xs []int, last int) bool {
			predicate := immutable_func.Konst[int](false)
			nonemptySlice := append(xs, last)
			xl := FromSlice(nonemptySlice)
			return Cycle(xl).Filter(predicate) == nil
		},
		"xs.filter(konst(true)) == xs": func(xs []int) bool {
			predicate := immutable_func.Konst[int](true)
			xl := FromSlice(xs)
			return slicesEqual(xl.Filter(predicate).ToSlice(), xs)
		},
		"cycle(xs).filter(konst(true)) == cycle(xs)": func(xs []int, last int) bool {
			predicate := immutable_func.Konst[int](true)
			nonemptySlice := append(xs, last)
			xl := FromSlice(nonemptySlice)
			return slicesEqual(Cycle(xl).Filter(predicate).Take(2*len(nonemptySlice)).ToSlice(), append(nonemptySlice, nonemptySlice...))
		},
	})
}

func TestListSort(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"sort(xs) is sorted": func(xs []int) bool {

			sortedXs := make([]int, len(xs))
			copy(sortedXs, xs)
			sort.Sort(sort.IntSlice(sortedXs))

			xl := FromSlice(xs)
			return slicesEqual(xl.Sort(comparator.OrderedComparator[int]).ToSlice(), sortedXs)
		},
		"sort(xs ++ cycle(ys)) == sort([x | x <- xs, x < min(ys)]) ++ repeat(min(ys))": func(xs []int, ys []int, last int) bool {
			nonemptySlice := append(ys, last)

			minY := math.MaxInt
			for _, y := range nonemptySlice {
				if minY > y {
					minY = y
				}
			}

			xl := FromSlice(xs)
			yl := FromSlice(nonemptySlice)
			cmp := comparator.OrderedComparator[int]
			result := xl.Append(Cycle(yl)).Sort(cmp)
			return result.IsIsomorphicTo(xl.Filter(func(x int) bool { return x < minY }).Sort(cmp).Append(Repeat(minY)), cmp)
		},
	})
}

func TestListReverse(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.reverse().length() == xs.length()": func(xs []int) bool {
			xl := FromSlice(xs)
			return xl.Reverse().Length() == len(xs)
		},
		"xs.reverse().reverse() == xs": func(xs []int) bool {
			xl := FromSlice(xs)
			return slicesEqual(xl.Reverse().Reverse().ToSlice(), xs)
		},
		"cons(x, xs).reverse() == xs.reverse().append([x])": func(xs []int, x int) bool {
			xl := FromSlice(xs)
			return slicesEqual(Cons(x, xl).Reverse().ToSlice(), append(xl.Reverse().ToSlice(), x))
		},
	})
}

func TestListIntersperse(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"[].intersperse(sep) == []": func(sep int) bool {
			return (*List[int])(nil).Intersperse(sep) == nil
		},
		"[x].intersperse(sep) == [x]": func(x int, sep int) bool {
			return slicesEqual(FromSlice([]int{x}).Intersperse(sep).ToSlice(), []int{x})
		},
		"cons(x, xs).intersperse(sep) == [x, sep] ++ xs.intersperse(sep)": func(xs []int, last int, x int, sep int) bool {
			nonemptySlice := append(xs, last)
			xl := FromSlice(nonemptySlice)
			return slicesEqual(Cons(x, xl).Intersperse(sep).ToSlice(), FromSlice([]int{x, sep}).Append(xl.Intersperse(sep)).ToSlice())
		},
		"repeat(x).intersperse(sep) == cycle([x, sep])": func(x int, sep int) bool {
			return Repeat(x).Intersperse(sep).IsIsomorphicTo(Cycle(FromSlice([]int{x, sep})), comparator.OrderedComparator[int])
		},
		"cycle(cons(x, xs)).intersperse(sep) == cycle([x, sep] ++ xs.intersperse(sep) ++ [sep])": func(xs []int, last int, x int, sep int) bool {
			nonemptySlice := append(xs, last)
			xl := FromSlice(nonemptySlice)
			return Cycle(Cons(x, xl)).Intersperse(sep).IsIsomorphicTo(Cycle(FromSlice([]int{x, sep}).Append(xl.Intersperse(sep)).Append(FromSlice([]int{sep}))), comparator.OrderedComparator[int])
		},
		"xs.append(cycle(ys)).intersperse(sep) == xs.intersperse(sep) ++ cons(sep, cycle(ys).intersperse(sep))": func(xs []int, ys []int, xLast int, yLast int, sep int) bool {
			nonemptySliceX := append(xs, xLast)
			nonemptySliceY := append(ys, yLast)
			xl := FromSlice(nonemptySliceX)
			yl := FromSlice(nonemptySliceY)
			return xl.Append(Cycle(yl)).Intersperse(sep).IsIsomorphicTo(xl.Intersperse(sep).Append(Cons(sep, Cycle(yl).Intersperse(sep))), comparator.OrderedComparator[int])
		},
	})
}

func TestListIsIsomorphicTo(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs is isomorphic to itself": func(xs []int) bool {
			xl := FromSlice(xs)
			return xl.IsIsomorphicTo(xl, comparator.OrderedComparator[int])
		},
		"xs.append(cycle(ys)) is isomorphic to xs.append(ys).append(cycle(ys))": func(xs []int, ys []int, last int) bool {
			nonemptySlice := append(ys, last)
			xl := FromSlice(xs)
			yl := FromSlice(nonemptySlice)
			return xl.Append(Cycle(yl)).
				IsIsomorphicTo(xl.Append(yl).Append(Cycle(yl)), comparator.OrderedComparator[int])
		},
	})
}

func TestListAll(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.append(ys).all(p) == xs.all(p) and ys.all(p)": func(xs []int, ys []int) bool {
			predicate := func(x int) bool { return x%100 < 90 }
			xl := FromSlice(xs)
			yl := FromSlice(ys)
			return xl.Append(yl).All(predicate) == (xl.All(predicate) && yl.All(predicate))
		},
		"cycle(xs).all(p) == xs.all(p)": func(xs []int, last int) bool {
			predicate := func(x int) bool { return x%100 < 90 }
			nonemptySlice := append(xs, last)
			xl := FromSlice(nonemptySlice)
			return Cycle(xl).All(predicate) == xl.All(predicate)
		},
		"repeat(x).all(p) == p(x)": func(x int) bool {
			predicate := func(x int) bool { return x%100 < 90 }
			return Repeat(x).All(predicate) == predicate(x)
		},
	})
}

func TestListAny(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.append(ys).any(p) == xs.any(p) or ys.any(p)": func(xs []int, ys []int) bool {
			predicate := func(x int) bool { return x%100 < 90 }
			xl := FromSlice(xs)
			yl := FromSlice(ys)
			return xl.Append(yl).Any(predicate) == (xl.Any(predicate) || yl.Any(predicate))
		},
		"cycle(xs).any(p) == xs.any(p)": func(xs []int, last int) bool {
			predicate := func(x int) bool { return x%100 < 90 }
			nonemptySlice := append(xs, last)
			xl := FromSlice(nonemptySlice)
			return Cycle(xl).Any(predicate) == xl.Any(predicate)
		},
		"repeat(x).any(p) == p(x)": func(x int) bool {
			predicate := func(x int) bool { return x%100 < 90 }
			return Repeat(x).Any(predicate) == predicate(x)
		},
	})
}

func checkProperties(t *testing.T, properties map[string]any) {
	for name, property := range properties {
		name, property := name, property
		t.Run(name, func(t *testing.T) {

			t.Parallel()

			err := quick.Check(property, nil)
			require.NoError(t, err)
		})
	}
}

func slicesEqual[T any](v1 []T, v2 []T) bool {
	return (len(v1) == 0 && len(v2) == 0) || reflect.DeepEqual(v1, v2)
}
