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
			xl := FromGoSlice(nonemptySlice)
			return !Cycle(xl).isFinite()
		},
		"cycle(xs).take(2 * xs.length()) == xs ++ xs": func(xs []int, last int) bool {
			nonemptySlice := append(xs, last)
			xl := FromGoSlice(nonemptySlice)
			return slicesEqual(Cycle(xl).Take(2*len(nonemptySlice)).ToGoSlice(), append(nonemptySlice, nonemptySlice...))
		},
		"cycle of an infinite list will be the same": func(prefixes []int, xs []int, last int) bool {
			nonemptySlice := append(xs, last)
			xl := Cycle(FromGoSlice(nonemptySlice))
			sampleLen := len(prefixes) + 2*len(nonemptySlice)
			return slicesEqual(Cycle(xl).Take(sampleLen).ToGoSlice(), xl.Take(sampleLen).ToGoSlice())
		},
	})
}

func TestMap(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.map(f).length() == xs.length()": func(xs []int) bool {
			f := func(x int) int { return x + 1 }
			l := FromGoSlice(xs)
			return Map(l, f).Length() == l.Length()
		},
		"cycle(xs).map(f) == cycle(xs.map(f))": func(xs []int, last int) bool {
			f := func(x int) int { return x + 1 }
			nonemptySlice := append(xs, last)
			xl := FromGoSlice(nonemptySlice)
			sampleLen := 2 * (len(xs) + 1)
			return slicesEqual(Map(Cycle(xl), f).Take(sampleLen).ToGoSlice(), Cycle(Map(xl, f)).Take(sampleLen).ToGoSlice())
		},
		"xs.map(f1).map(f2) == xs.map(f2 . f1)": func(xs []int) bool {
			f1 := func(x int) int { return x + 1 }
			f2 := strconv.Itoa
			xl := FromGoSlice(xs)
			return slicesEqual(Map(Map(xl, f1), f2).ToGoSlice(), Map(xl, func(x int) string { return f2(f1(x)) }).ToGoSlice())
		},
	})
}

func TestFoldl(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.foldl([], flip(cons)).reverse() == xs": func(xs []int) bool {
			xl := FromGoSlice(xs)
			return slicesEqual(Foldl(xl, nil, func(acc *List[int], x int) *List[int] { return Cons(x, acc) }).Reverse().ToGoSlice(), xs)
		},
	})
}

func TestFoldr(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.foldr([], cons) == xs": func(xs []int) bool {
			xl := FromGoSlice(xs)
			return slicesEqual(Foldr(xl, nil, Cons[int]).ToGoSlice(), xs)
		},
	})
}

func TestConcat(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"concat([[]] * N) == []": func(n uint) bool {
			n %= 100
			return Concat(Repeat((*List[int])(nil))).Take(int(n)) == nil
		},
		"concat(xss) == xss.foldl([], ++)": func(xss [][]int) bool {
			xll := Map(FromGoSlice(xss), FromGoSlice[int])
			return slicesEqual(Concat(xll).ToGoSlice(), Foldl(xll, (*List[int])(nil), (*List[int]).Append).ToGoSlice())
		},
		"concat(repeat(xs)) == cycle(xs)": func(xs []int, last int) bool {
			nonemptySlice := append(xs, last)
			xl := FromGoSlice(nonemptySlice)
			return Concat(Repeat(xl)).IsIsomorphicTo(Cycle(xl), comparator.OrderedComparator[int])
		},
		"concat([cycle(xs)]) == cycle(xs)": func(xs []int, last int) bool {
			nonemptySlice := append(xs, last)
			xl := FromGoSlice(nonemptySlice)
			return Concat(Cons(Cycle(xl), nil)).IsIsomorphicTo(Cycle(xl), comparator.OrderedComparator[int])
		},
	})
}

func TestListAppend(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.append(ys) == xs ++ ys": func(xs []int, ys []int) bool {
			xl := FromGoSlice(xs)
			yl := FromGoSlice(ys)
			return slicesEqual(xl.Append(yl).ToGoSlice(), append(xs, ys...))
		},
		"cycle(xs).append(ys) == cycle(xs)": func(xs []int, ys []int, xLast int) bool {
			nonemptySliceX := append(xs, xLast)
			xl := FromGoSlice(nonemptySliceX)
			yl := FromGoSlice(ys)
			return Cycle(xl).Append(yl).IsIsomorphicTo(Cycle(xl), comparator.OrderedComparator[int])
		},
	})
}

func TestListTake(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.append(ys).take(xs.length()) == xs": func(xs []int, ys []int) bool {
			xl := FromGoSlice(xs)
			yl := FromGoSlice(ys)
			return slicesEqual(xl.Append(yl).Take(len(xs)).ToGoSlice(), xs)
		},
		"xs.take(n) == xs if n >= xs.length()": func(xs []int, delta uint8) bool {
			n := len(xs) + int(delta)
			xl := FromGoSlice(xs)
			return slicesEqual(xl.Take(n).ToGoSlice(), xs)
		},
	})
}

func TestListDrop(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.append(ys).drop(xs.length()) == ys": func(xs []int, ys []int) bool {
			xl := FromGoSlice(xs)
			yl := FromGoSlice(ys)
			return slicesEqual(xl.Append(yl).Drop(len(xs)).ToGoSlice(), ys)
		},
		"xs.drop(n) == nil if n >= xs.length()": func(xs []int, delta uint8) bool {
			n := len(xs) + int(delta)
			xl := FromGoSlice(xs)
			return xl.Drop(n) == nil
		},
	})
}

func TestListFilter(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.filter(p).append(ys.filter(p)) == xs.append(ys).filter(p)": func(xs []int, ys []int) bool {
			predicate := func(x int) bool { return x%2 == 0 }
			xl := FromGoSlice(xs)
			yl := FromGoSlice(ys)
			return reflect.DeepEqual(xl.Filter(predicate).Append(yl.Filter(predicate)).ToGoSlice(), xl.Append(yl).Filter(predicate).ToGoSlice())
		},
		"xs.filter(konst(false)) == nil": func(xs []int) bool {
			predicate := immutable_func.Konst[int](false)
			xl := FromGoSlice(xs)
			return xl.Filter(predicate) == nil
		},
		"cycle(xs).filter(konst(false)) == nil": func(xs []int, last int) bool {
			predicate := immutable_func.Konst[int](false)
			nonemptySlice := append(xs, last)
			xl := FromGoSlice(nonemptySlice)
			return Cycle(xl).Filter(predicate) == nil
		},
		"xs.filter(konst(true)) == xs": func(xs []int) bool {
			predicate := immutable_func.Konst[int](true)
			xl := FromGoSlice(xs)
			return slicesEqual(xl.Filter(predicate).ToGoSlice(), xs)
		},
		"cycle(xs).filter(konst(true)) == cycle(xs)": func(xs []int, last int) bool {
			predicate := immutable_func.Konst[int](true)
			nonemptySlice := append(xs, last)
			xl := FromGoSlice(nonemptySlice)
			return slicesEqual(Cycle(xl).Filter(predicate).Take(2*len(nonemptySlice)).ToGoSlice(), append(nonemptySlice, nonemptySlice...))
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

			xl := FromGoSlice(xs)
			return slicesEqual(xl.Sort(comparator.OrderedComparator[int]).ToGoSlice(), sortedXs)
		},
		"sort(xs ++ cycle(ys)) == sort([x | x <- xs, x < min(ys)]) ++ repeat(min(ys))": func(xs []int, ys []int, last int) bool {
			nonemptySlice := append(ys, last)

			minY := math.MaxInt
			for _, y := range nonemptySlice {
				if minY > y {
					minY = y
				}
			}

			xl := FromGoSlice(xs)
			yl := FromGoSlice(nonemptySlice)
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
			xl := FromGoSlice(xs)
			return xl.Reverse().Length() == len(xs)
		},
		"xs.reverse().reverse() == xs": func(xs []int) bool {
			xl := FromGoSlice(xs)
			return slicesEqual(xl.Reverse().Reverse().ToGoSlice(), xs)
		},
		"cons(x, xs).reverse() == xs.reverse().append([x])": func(xs []int, x int) bool {
			xl := FromGoSlice(xs)
			return slicesEqual(Cons(x, xl).Reverse().ToGoSlice(), append(xl.Reverse().ToGoSlice(), x))
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
			return slicesEqual(FromGoSlice([]int{x}).Intersperse(sep).ToGoSlice(), []int{x})
		},
		"cons(x, xs).intersperse(sep) == [x, sep] ++ xs.intersperse(sep)": func(xs []int, last int, x int, sep int) bool {
			nonemptySlice := append(xs, last)
			xl := FromGoSlice(nonemptySlice)
			return slicesEqual(Cons(x, xl).Intersperse(sep).ToGoSlice(), FromGoSlice([]int{x, sep}).Append(xl.Intersperse(sep)).ToGoSlice())
		},
		"repeat(x).intersperse(sep) == cycle([x, sep])": func(x int, sep int) bool {
			return Repeat(x).Intersperse(sep).IsIsomorphicTo(Cycle(FromGoSlice([]int{x, sep})), comparator.OrderedComparator[int])
		},
		"cycle(cons(x, xs)).intersperse(sep) == cycle([x, sep] ++ xs.intersperse(sep) ++ [sep])": func(xs []int, last int, x int, sep int) bool {
			nonemptySlice := append(xs, last)
			xl := FromGoSlice(nonemptySlice)
			return Cycle(Cons(x, xl)).Intersperse(sep).IsIsomorphicTo(Cycle(FromGoSlice([]int{x, sep}).Append(xl.Intersperse(sep)).Append(FromGoSlice([]int{sep}))), comparator.OrderedComparator[int])
		},
		"xs.append(cycle(ys)).intersperse(sep) == xs.intersperse(sep) ++ cons(sep, cycle(ys).intersperse(sep))": func(xs []int, ys []int, xLast int, yLast int, sep int) bool {
			nonemptySliceX := append(xs, xLast)
			nonemptySliceY := append(ys, yLast)
			xl := FromGoSlice(nonemptySliceX)
			yl := FromGoSlice(nonemptySliceY)
			return xl.Append(Cycle(yl)).Intersperse(sep).IsIsomorphicTo(xl.Intersperse(sep).Append(Cons(sep, Cycle(yl).Intersperse(sep))), comparator.OrderedComparator[int])
		},
	})
}

func TestListIsIsomorphicTo(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs is isomorphic to itself": func(xs []int) bool {
			xl := FromGoSlice(xs)
			return xl.IsIsomorphicTo(xl, comparator.OrderedComparator[int])
		},
		"xs.append(cycle(ys)) is isomorphic to xs.append(ys).append(cycle(ys))": func(xs []int, ys []int, last int) bool {
			nonemptySlice := append(ys, last)
			xl := FromGoSlice(xs)
			yl := FromGoSlice(nonemptySlice)
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
			xl := FromGoSlice(xs)
			yl := FromGoSlice(ys)
			return xl.Append(yl).All(predicate) == (xl.All(predicate) && yl.All(predicate))
		},
		"cycle(xs).all(p) == xs.all(p)": func(xs []int, last int) bool {
			predicate := func(x int) bool { return x%100 < 90 }
			nonemptySlice := append(xs, last)
			xl := FromGoSlice(nonemptySlice)
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
			xl := FromGoSlice(xs)
			yl := FromGoSlice(ys)
			return xl.Append(yl).Any(predicate) == (xl.Any(predicate) || yl.Any(predicate))
		},
		"cycle(xs).any(p) == xs.any(p)": func(xs []int, last int) bool {
			predicate := func(x int) bool { return x%100 < 90 }
			nonemptySlice := append(xs, last)
			xl := FromGoSlice(nonemptySlice)
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
