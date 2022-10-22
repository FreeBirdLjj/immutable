package immutable_set

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetInsert(t *testing.T) {

	t.Parallel()

	t.Run("should succeed", func(t *testing.T) {

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

		s := NewSet(strings.Compare)
		for _, value := range values {
			s, _ = s.Insert(value)
		}

		for _, value := range values {
			assert.True(t, s.Has(value))
		}
	})
}

func TestSetDelete(t *testing.T) {

	t.Parallel()

	t.Run("should succeed", func(t *testing.T) {

		t.Parallel()

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		values := []string{
			"one",
			"two",
			"three",
			"four",
			"five",
			"six",
			"seven",
			"eight",
		}

		r.Shuffle(len(values), func(i int, j int) {
			values[i], values[j] = values[j], values[i]
		})

		s := NewSetFromValues(strings.Compare, values...)

		r.Shuffle(len(values), func(i int, j int) {
			values[i], values[j] = values[j], values[i]
		})

		for i, value := range values {

			s, _ = s.Delete(value)

			for _, deletedValue := range values[:i+1] {
				assert.False(t, s.Has(deletedValue))
			}

			for _, remainingValue := range values[i+1:] {
				assert.True(t, s.Has(remainingValue))
			}
		}
	})
}
