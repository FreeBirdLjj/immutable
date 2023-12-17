package slice

import (
	"reflect"
	"sort"
	"strconv"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/require"

	"github.com/freebirdljj/immutable/comparator"
	immutable_func "github.com/freebirdljj/immutable/func"
)

func TestMap(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.map(f).length() == len(xs)": func(xs []int) bool {
			f := func(x int) int { return x + 1 }
			l := FromGoSlice(xs)
			return len(Map(l, f)) == len(l)
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
		"xs.foldl([], append) == xs": func(xs []int) bool {
			xl := FromGoSlice(xs)
			return slicesEqual(Foldl(xl, nil, func(acc Slice[int], x int) Slice[int] { return append(acc, x) }).ToGoSlice(), xs)
		},
	})
}

func TestFoldr(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.foldr([], append).reverse() == xs": func(xs []int) bool {
			xl := FromGoSlice(xs)
			return slicesEqual(Foldr(xl, nil, func(x int, acc Slice[int]) Slice[int] { return append(acc, x) }).Reverse().ToGoSlice(), xs)
		},
	})
}

func TestMaximumBy(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"`maximumBy()` returns the max": func(xs []int, lastX int) bool {
			nonemptySlice := append(xs, lastX)

			max := lastX
			for _, x := range xs {
				if max < x {
					max = x
				}
			}

			xl := FromGoSlice(nonemptySlice)
			return MaximumBy(xl, comparator.OrderedComparator[int]) == max
		},
	})
}

func TestMinimumBy(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"`mainimumBy()` returns the min": func(xs []int, lastX int) bool {
			nonemptySlice := append(xs, lastX)

			min := lastX
			for _, x := range xs {
				if min > x {
					min = x
				}
			}

			xl := FromGoSlice(nonemptySlice)
			return MinimumBy(xl, comparator.OrderedComparator[int]) == min
		},
	})
}

func TestConcat(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"concat([[]] * N) == []": func(n uint) bool {
			n %= 100
			return len(Concat(make(Slice[Slice[int]], n))) == 0
		},
		"concat(xss) == xss.foldl([], ++)": func(xss [][]int) bool {
			xll := Map(FromGoSlice(xss), FromGoSlice[int])
			return slicesEqual(Concat(xll).ToGoSlice(), Foldl(xll, Slice[int](nil), func(acc Slice[int], xs Slice[int]) Slice[int] { return acc.Append(xs...) }).ToGoSlice())
		},
	})
}

func TestSliceAppend(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.append(elems...) == append(xs, elems...)": func(xs []int, elems []int) bool {
			xl := FromGoSlice(xs)
			return slicesEqual(xl.Append(elems...).ToGoSlice(), append(xs, elems...))
		},
	})
}

func TestSliceTake(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.append(ys).take(len(xs)) == xs": func(xs []int, ys []int) bool {
			xl := FromGoSlice(xs)
			return slicesEqual(xl.Append(ys...).Take(len(xs)).ToGoSlice(), xs)
		},
	})
}

func TestSliceDrop(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.append(ys).drop(xs.length()) == ys": func(xs []int, ys []int) bool {
			xl := FromGoSlice(xs)
			return slicesEqual(xl.Append(ys...).Drop(len(xs)).ToGoSlice(), ys)
		},
	})
}

func TestSliceFind(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.find(konst(false)) == nil": func(xs []int) bool {
			predicate := immutable_func.Konst[int](false)
			xl := FromGoSlice(xs)
			return xl.Find(predicate) == nil
		},
		"p(x) == true -> *([x].append(xs).find(p) == x": func(x int, xs []int) bool {
			predicate := func(val int) bool { return val == x }
			xl := FromGoSlice([]int{x}).Append(xs...)
			return reflect.DeepEqual(xl.Find(predicate), &x)
		},
		"p(x) == false -> [x].append(xs).find(p) == xs.find(p)": func(x int, xs []int) bool {
			predicate := func(val int) bool { return val%2 != x%2 }
			return reflect.DeepEqual(
				FromGoSlice([]int{x}).Append(xs...).Find(predicate),
				FromGoSlice(xs).Find(predicate),
			)
		},
	})
}

func TestSliceFilter(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.filter(p).append(ys.filter(p)) == xs.append(ys).filter(p)": func(xs []int, ys []int) bool {
			predicate := func(x int) bool { return x%2 == 0 }
			xl := FromGoSlice(xs)
			yl := FromGoSlice(ys)
			return reflect.DeepEqual(xl.Filter(predicate).Append(yl.Filter(predicate)...).ToGoSlice(), xl.Append(yl...).Filter(predicate).ToGoSlice())
		},
		"xs.filter(konst(false)) == []": func(xs []int) bool {
			predicate := immutable_func.Konst[int](false)
			xl := FromGoSlice(xs)
			return len(xl.Filter(predicate)) == 0
		},
		"xs.filter(konst(true)) == xs": func(xs []int) bool {
			predicate := immutable_func.Konst[int](true)
			xl := FromGoSlice(xs)
			return slicesEqual(xl.Filter(predicate).ToGoSlice(), xs)
		},
	})
}

func TestSliceSort(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"sort(xs) is sorted": func(xs []int) bool {

			sortedXs := make([]int, len(xs))
			copy(sortedXs, xs)
			sort.Sort(sort.IntSlice(sortedXs))

			xl := FromGoSlice(xs)
			return slicesEqual(xl.Sort(comparator.OrderedComparator[int]).ToGoSlice(), sortedXs)
		},
	})
}

func TestSliceReverse(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.reverse().length() == xs.length()": func(xs []int) bool {
			xl := FromGoSlice(xs)
			return len(xl.Reverse()) == len(xs)
		},
		"xs.reverse().reverse() == xs": func(xs []int) bool {
			xl := FromGoSlice(xs)
			return slicesEqual(xl.Reverse().Reverse().ToGoSlice(), xs)
		},
		"append(xs, x).reverse().tail() == xs.reverse()": func(xs []int, x int) bool {
			xl := FromGoSlice(xs)
			return slicesEqual(xl.Append(x).Reverse().Tail().ToGoSlice(), xl.Reverse().ToGoSlice())
		},
	})
}

func TestSliceIntersperse(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"[].intersperse(sep) == []": func(sep int) bool {
			return Slice[int](nil).Intersperse(sep) == nil
		},
		"[x].intersperse(sep) == [x]": func(x int, sep int) bool {
			return slicesEqual(FromGoSlice([]int{x}).Intersperse(sep).ToGoSlice(), []int{x})
		},
		"xs.append(x).intersperse(sep) == xs.intersperse(sep).append(sep, x)": func(xs []int, last int, x int, sep int) bool {
			nonemptySlice := append(xs, last)
			xl := FromGoSlice(nonemptySlice)
			return slicesEqual(xl.Append(x).Intersperse(sep).ToGoSlice(), xl.Intersperse(sep).Append(sep, x).ToGoSlice())
		},
	})
}

func TestSliceAll(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.append(ys).all(p) == xs.all(p) and ys.all(p)": func(xs []int, ys []int) bool {
			predicate := func(x int) bool { return x%100 < 90 }
			xl := FromGoSlice(xs)
			yl := FromGoSlice(ys)
			return xl.Append(yl...).All(predicate) == (xl.All(predicate) && yl.All(predicate))
		},
		"xs.all(konst(true)) == true": func(xs []int) bool {
			predicate := immutable_func.Konst[int](true)
			xl := FromGoSlice(xs)
			return xl.All(predicate)
		},
	})
}

func TestSliceAny(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"xs.append(ys).any(p) == xs.any(p) or ys.any(p)": func(xs []int, ys []int) bool {
			predicate := func(x int) bool { return x%100 < 90 }
			xl := FromGoSlice(xs)
			yl := FromGoSlice(ys)
			return xl.Append(yl...).Any(predicate) == (xl.Any(predicate) || yl.Any(predicate))
		},
		"xs.any(konst(false)) == false": func(xs []int) bool {
			predicate := immutable_func.Konst[int](false)
			xl := FromGoSlice(xs)
			return !xl.Any(predicate)
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
