package iter

import (
	"iter"
	"slices"
	"strconv"
	"testing"

	"github.com/freebirdljj/immutable/internal/quick"
)

func TestEmpty(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"Empty() == []": func() bool {
			return slices.Equal(
				slices.Collect(Empty[int]()),
				nil,
			)
		},
	})
}

func TestSingleton(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"Singoleton(x) == [x]": func(x int) bool {
			return slices.Equal(
				slices.Collect(Singleton(x)),
				[]int{x},
			)
		},
	})
}

func TestCycle(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"Take(Cycle(xs), 2 * len(xs)) == xs ++ xs": func(xs []int, last int) bool {
			nonemptySlice := append(xs, last)
			return slices.Equal(
				slices.Collect(Take(Cycle(slices.Values(nonemptySlice)), 2*len(nonemptySlice))),
				append(nonemptySlice, nonemptySlice...),
			)
		},
		"cycle of an infinite list will be the same": func(prefixes []int, xs []int, last int) bool {
			nonemptySlice := append(xs, last)
			seq := Cycle(slices.Values(nonemptySlice))
			sampleLen := len(prefixes) + 2*len(nonemptySlice)
			return slices.Equal(
				slices.Collect(Take(Cycle(seq), sampleLen)),
				slices.Collect(Take(seq, sampleLen)),
			)
		},
	})
}

func TestMap(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"Map(Cycle(seq, f)) == Cycle(Map(seq, f))": func(xs []int, last int) bool {
			f := func(x int) int { return x + 1 }
			nonemptySlice := append(xs, last)
			seq := slices.Values(nonemptySlice)
			sampleLen := 2 * (len(xs) + 1)
			return slices.Equal(
				slices.Collect(Take(Map(Cycle(seq), f), sampleLen)),
				slices.Collect(Take(Cycle(Map(seq, f)), sampleLen)),
			)
		},
		"Map(Map(seq, f1), f2) == Map(seq, f1 . f2)": func(xs []int) bool {
			f1 := func(x int) int { return x + 1 }
			f2 := strconv.Itoa
			seq := slices.Values(xs)
			return slices.Equal(
				slices.Collect(Map(Map(seq, f1), f2)),
				slices.Collect(Map(seq, func(x int) string { return f2(f1(x)) })),
			)
		},
	})
}

func TestConcat(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"Concat([[]] * N) == []": func(n uint) bool {
			n %= 100
			seq := func(yield func(iter.Seq[int]) bool) {
				for i := 0; i < int(n); i++ {
					if !yield(Empty[int]()) {
						return
					}
				}
			}
			return slices.Equal(
				slices.Collect(Concat(seq)),
				nil,
			)
		},
		"Concat([xs, ys]) == xs ++ ys": func(xs []int, ys []int) bool {
			return slices.Equal(
				slices.Collect(Concat(slices.Values([]iter.Seq[int]{
					slices.Values(xs),
					slices.Values(ys),
				}))),
				append(xs, ys...),
			)
		},
	})
}

func TestTake(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"Take(xs ++ ys, len(xs)) == xs": func(xs []int, ys []int) bool {
			return slices.Equal(
				slices.Collect(Take(slices.Values(append(xs, ys...)), len(xs))),
				xs,
			)
		},
		"Take(xs, n) == xs if n >= len(xs)": func(xs []int, delta uint8) bool {
			n := len(xs) + int(delta)
			seq := slices.Values(xs)
			return slices.Equal(
				slices.Collect(Take(seq, n)),
				xs,
			)
		},
	})
}

func TestDrop(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"Drop(xs ++ ys, len(xs)) == ys": func(xs []int, ys []int) bool {
			return slices.Equal(
				slices.Collect(Drop(slices.Values(append(xs, ys...)), len(xs))),
				ys,
			)
		},
		"Drop(xs, n) == [] if n >= xs.length()": func(xs []int, delta uint8) bool {
			n := len(xs) + int(delta)
			return slices.Equal(
				slices.Collect(Drop(slices.Values(xs), n)),
				nil,
			)
		},
	})
}
