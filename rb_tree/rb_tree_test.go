package immutable_rb_tree

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/freebirdljj/immutable/comparator"
	"github.com/freebirdljj/immutable/internal/quick"
)

func TestRBTreeInsert(t *testing.T) {

	t.Parallel()

	t.Run("should succeed to insert new node", func(t *testing.T) {

		t.Parallel()

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		values := []string{
			"one",
			"two",
			"three",
			"four",
			"five",
			"six",
		}

		r.Shuffle(len(values), func(i int, j int) {
			values[i], values[j] = values[j], values[i]
		})

		rbTree := New(strings.Compare)
		for _, value := range values {

			newRBTree, affected := rbTree.Insert(value)
			assert.True(t, affected)

			rbTree = newRBTree
		}

		gotValues := rbTree.Values()
		assert.Len(t, values, len(gotValues))
		assert.Equal(t, len(values), rbTree.Count())

		for _, value := range values {
			assert.Contains(t, gotValues, value)
		}
	})
	t.Run("should succeed to update existing node", func(t *testing.T) {

		t.Parallel()

		N := 10

		rbTree := New(func(l int, r int) int { return (l % N) - (r % N) })

		for i := 0; i < N; i++ {
			rbTree, _ = rbTree.Insert(i)
		}

		for i := N; i < N*2; i++ {

			newRBTree, affected := rbTree.Insert(i)
			assert.False(t, affected)

			rbTree = newRBTree
		}

		gotValues := rbTree.Values()
		assert.Len(t, gotValues, N)

		for i := N; i < N*2; i++ {
			assert.Contains(t, gotValues, i)
		}
	})
}

func TestRBTreeDelete(t *testing.T) {

	t.Parallel()

	t.Run("should succeed to delete existing node", func(t *testing.T) {

		t.Parallel()

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		values := []string{
			"one",
			"two",
			"three",
			"four",
			"five",
			"six",
		}

		r.Shuffle(len(values), func(i int, j int) {
			values[i], values[j] = values[j], values[i]
		})

		rbTree := FromValues(strings.Compare, values...)

		r.Shuffle(len(values), func(i int, j int) {
			values[i], values[j] = values[j], values[i]
		})

		for _, value := range values {

			newRBTree, affected := rbTree.Delete(value)
			assert.True(t, affected)
			assert.Nil(t, newRBTree.Lookup(value))

			rbTree = newRBTree
		}

		assert.Empty(t, rbTree.Values())
		assert.Zero(t, rbTree.Count())
	})
	t.Run("should succeed to delete non-existing node", func(t *testing.T) {

		t.Parallel()

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		values := []string{
			"one",
			"two",
			"three",
			"four",
			"five",
			"six",
		}
		nonExistingValue := "zero"

		r.Shuffle(len(values), func(i int, j int) {
			values[i], values[j] = values[j], values[i]
		})

		rbTree := FromValues(strings.Compare, values...)

		newRBTree, affected := rbTree.Delete(nonExistingValue)
		assert.False(t, affected)

		assert.Len(t, newRBTree.Values(), len(values))
		assert.Equal(t, len(values), newRBTree.Count())
	})
}

func TestRBTreeMaximum(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"rb_tree.fromValues(xs).maximum() == max(xs)": func(xs []int, lastX int) bool {

			nonemptySlice := append(xs, lastX)

			max := lastX
			for _, x := range xs {
				if max < x {
					max = x
				}
			}

			rbTree := FromValues(comparator.OrderedComparator[int], nonemptySlice...)
			return rbTree.Maximum() == max
		},
	})
}

func TestRBTreeMinimum(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"rb_tree.fromValues(xs).minimum() == min(xs)": func(xs []int, lastX int) bool {

			nonemptySlice := append(xs, lastX)

			min := lastX
			for _, x := range xs {
				if min > x {
					min = x
				}
			}

			rbTree := FromValues(comparator.OrderedComparator[int], nonemptySlice...)
			return rbTree.Minimum() == min
		},
	})
}
