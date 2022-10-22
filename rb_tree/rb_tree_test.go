package immutable_rb_tree

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

		rbTree := NewRBTree(strings.Compare)
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

		rbTree := NewRBTree(func(l int, r int) int { return (l % N) - (r % N) })

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

		rbTree := NewRBTreeFromValues(strings.Compare, values...)

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

		rbTree := NewRBTreeFromValues(strings.Compare, values...)

		newRBTree, affected := rbTree.Delete(nonExistingValue)
		assert.False(t, affected)

		assert.Len(t, newRBTree.Values(), len(values))
		assert.Equal(t, len(values), newRBTree.Count())
	})
}
