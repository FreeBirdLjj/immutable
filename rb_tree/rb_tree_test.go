package immutable_rb_tree

import (
	"reflect"
	"slices"
	"testing"

	"github.com/freebirdljj/immutable/comparator"
	"github.com/freebirdljj/immutable/internal/quick"
)

func TestRBTreeInsert(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"should succeed to insert a new node": func(xs []int, x int) bool {

			newXs := make([]int, 0, len(xs))
			for _, value := range xs {
				if value != x {
					newXs = append(newXs, value)
				}
			}

			rbTree := FromValues(comparator.OrderedComparator[int], newXs...)

			newRBTree, affected := rbTree.Insert(x)
			return affected && slices.Contains(newRBTree.Values(), x)
		},
		"should succeed to update an existing node": func(xs []int, x int) bool {

			const N = 10
			updatedValue := x + N

			rbTree := FromValues(
				comparator.CascadeComparator(
					comparator.OrderedComparator[int],
					func(value int) int { return value % N },
				),
				append(xs, x)...,
			)

			newRBTree, affected := rbTree.Insert(updatedValue)
			values := newRBTree.Values()
			return !affected &&
				slices.Contains(values, updatedValue) &&
				!slices.Contains(values, x)
		},
	})
}

func TestRBTreeDelete(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"should succeed to delete an existing node": func(xs []int, x int) bool {

			rbTree := FromValues(comparator.OrderedComparator[int], append(xs, x)...)

			newRBTree, affected := rbTree.Delete(x)
			return affected && !slices.Contains(newRBTree.Values(), x)
		},
		"should succeed to delete non-existing node": func(xs []int, x int) bool {

			rbTree := New(comparator.OrderedComparator[int])

			for _, value := range xs {
				if value != x {
					rbTree, _ = rbTree.Insert(value)
				}
			}

			newRBTree, affected := rbTree.Delete(x)
			return !affected && newRBTree == rbTree
		},
	})
}

func TestRBTreeMaximum(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"rb_tree.fromValues(xs).maximum() == max(xs)": func(xs []int, lastX int) bool {

			nonemptySlice := append(xs, lastX)
			max := slices.Max(nonemptySlice)

			rbTree := FromValues(comparator.OrderedComparator[int], nonemptySlice...)
			return rbTree.Maximum() == max
		},
	})
}

func TestRBTreeMinimum(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"rb_tree.fromValues(xs).minimum() == min(xs)": func(xs []int, lastX int) bool {

			nonemptySlice := append(xs, lastX)
			min := slices.Min(nonemptySlice)

			rbTree := FromValues(comparator.OrderedComparator[int], nonemptySlice...)
			return rbTree.Minimum() == min
		},
	})
}

func TestRBTreeInorderTraversal(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"rbTree.InorderTraversal(Konst(false)) should only iterate over at most 1 value": func(xs []int, lastX int) bool {

			nonemptySlice := append(xs, lastX)

			rbTree := FromValues(comparator.OrderedComparator[int], nonemptySlice...)

			cnt := 0
			rbTree.InorderTraversal()(func(value int) bool {
				cnt++
				return false
			})

			return cnt == 1
		},
		"rbTree.InorderTraversal(Konst(true)) should iterate over all values": func(xs []int, lastX int) bool {

			nonemptySlice := append(xs, lastX)

			allValues := make(map[int]struct{}, len(nonemptySlice))
			for _, value := range nonemptySlice {
				allValues[value] = struct{}{}
			}

			rbTree := FromValues(comparator.OrderedComparator[int], nonemptySlice...)

			visited := make(map[int]struct{}, len(nonemptySlice))
			rbTree.InorderTraversal()(func(value int) bool {
				visited[value] = struct{}{}
				return true
			})

			return reflect.DeepEqual(visited, allValues)
		},
	})
}
