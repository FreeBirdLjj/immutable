package list

import (
	"reflect"
	"strconv"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/require"

	immutable_func "github.com/freebirdljj/immutable/func"
)

func TestCycle(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]interface{}{
		"cycle(xs) is infinite": func(xs []int, last int) bool {
			nonemptySlice := append(xs, last)
			xl := FromSlice(nonemptySlice)
			return !Cycle(xl).isFinite()
		},
		"cycle(xs).take(2 * xs.length()) == xs ++ xs": func(xs []int, last int) bool {
			nonemptySlice := append(xs, last)
			xl := FromSlice(nonemptySlice)
			return reflect.DeepEqual(Cycle(xl).Take(2*len(nonemptySlice)).ToSlice(), append(nonemptySlice, nonemptySlice...))
		},
		"cycle of an infinite list will be the same": func(prefixes []int, xs []int, last int) bool {
			nonemptySlice := append(xs, last)
			xl := Cycle(FromSlice(nonemptySlice))
			sampleLen := len(prefixes) + 2*len(nonemptySlice)
			return reflect.DeepEqual(Cycle(xl).Take(sampleLen).ToSlice(), xl.Take(sampleLen).ToSlice())
		},
	})
}

func TestMap(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]interface{}{
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
			return reflect.DeepEqual(Map(Cycle(xl), f).Take(sampleLen).ToSlice(), Cycle(Map(xl, f)).Take(sampleLen).ToSlice())
		},
		"xs.map(f1).map(f2) == xs.map(f2 . f1)": func(xs []int) bool {
			f1 := func(x int) int { return x + 1 }
			f2 := strconv.Itoa
			xl := FromSlice(xs)
			return reflect.DeepEqual(Map(Map(xl, f1), f2).ToSlice(), Map(xl, func(x int) string { return f2(f1(x)) }).ToSlice())
		},
	})
}

func TestListAppend(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]interface{}{
		"xs.append(ys) == xs ++ ys": func(xs []int, ys []int) bool {
			xl := FromSlice(xs)
			yl := FromSlice(ys)
			return reflect.DeepEqual(xl.Append(yl).ToSlice(), append(xs, ys...))
		},
	})
}

func TestListTake(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]interface{}{
		"xs.append(ys).take(xs.length()) == xs": func(xs []int, ys []int) bool {
			if len(xs) == 0 {
				xs = nil
			}
			xl := FromSlice(xs)
			yl := FromSlice(ys)
			return reflect.DeepEqual(xl.Append(yl).Take(len(xs)).ToSlice(), xs)
		},
		"xs.take(n) == xs if n >= xs.length()": func(xs []int, delta uint8) bool {
			n := len(xs) + int(delta)
			if len(xs) == 0 {
				xs = nil
			}
			xl := FromSlice(xs)
			return reflect.DeepEqual(xl.Take(n).ToSlice(), xs)
		},
	})
}

func TestListDrop(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]interface{}{
		"xs.append(ys).drop(xs.length()) == ys": func(xs []int, ys []int) bool {
			if len(ys) == 0 {
				ys = nil
			}
			xl := FromSlice(xs)
			yl := FromSlice(ys)
			return reflect.DeepEqual(xl.Append(yl).Drop(len(xs)).ToSlice(), ys)
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

	checkProperties(t, map[string]interface{}{
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
			if len(xs) == 0 {
				xs = nil
			}
			xl := FromSlice(xs)
			return reflect.DeepEqual(xl.Filter(predicate).ToSlice(), xs)
		},
		"cycle(xs).filter(konst(true)) == cycle(xs)": func(xs []int, last int) bool {
			predicate := immutable_func.Konst[int](true)
			nonemptySlice := append(xs, last)
			xl := FromSlice(nonemptySlice)
			return reflect.DeepEqual(Cycle(xl).Filter(predicate).Take(2*len(nonemptySlice)).ToSlice(), append(nonemptySlice, nonemptySlice...))
		},
	})
}

func checkProperties(t *testing.T, properties map[string]interface{}) {
	for name, property := range properties {
		name, property := name, property
		t.Run(name, func(t *testing.T) {

			t.Parallel()

			err := quick.Check(property, nil)
			require.NoError(t, err)
		})
	}
}
