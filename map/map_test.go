package immutable_map

import (
	"slices"
	"strconv"
	"testing"

	"github.com/freebirdljj/immutable/comparator"
	"github.com/freebirdljj/immutable/internal/quick"
	"github.com/freebirdljj/immutable/tuple"
)

func TestMapInsert(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"should succeed to insert a new key": func(xs []int, x int) bool {

			formatter := func(value int) string { return strconv.FormatInt(int64(value), 10) }

			newXs := make([]int, 0, len(xs))
			for _, value := range xs {
				if value != x {
					newXs = append(newXs, value)
				}
			}

			m := New[int, string](comparator.OrderedComparator[int])

			for _, value := range newXs {
				m, _ = m.Insert(value, formatter(value))
			}

			xStr := formatter(x)
			newM, affected := m.Insert(x, xStr)
			return affected &&
				slices.Contains(newM.KeyValuePairs(), tuple.KeyValuePair[int, string]{Key: x, Value: xStr})
		},
		"should succeed to update an existing key": func(xs []int, x int) bool {

			formatter := func(value int) string { return strconv.FormatInt(int64(value), 10) }

			m := New[int, string](comparator.OrderedComparator[int])

			for _, value := range append(xs, x) {
				m, _ = m.Insert(value, formatter(value))
			}

			newValue := "new value"
			newM, affected := m.Insert(x, newValue)
			kvs := newM.KeyValuePairs()
			return !affected &&
				slices.Contains(kvs, tuple.KeyValuePair[int, string]{Key: x, Value: newValue}) &&
				!slices.Contains(kvs, tuple.KeyValuePair[int, string]{Key: x, Value: formatter(x)})
		},
	})
}

func TestMapDelete(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"should succeed to delete an existing key-value pair": func(xs []int, x int) bool {

			formatter := func(value int) string { return strconv.FormatInt(int64(value), 10) }

			m := New[int, string](comparator.OrderedComparator[int])

			for _, value := range append(xs, x) {
				m, _ = m.Insert(value, formatter(value))
			}

			newM, affected := m.Delete(x)
			return affected &&
				!slices.Contains(newM.KeyValuePairs(), tuple.KeyValuePair[int, string]{Key: x, Value: formatter(x)})
		},
		"should succeed to delete a non-existing key-value pair": func(xs []int, x int) bool {

			formatter := func(value int) string { return strconv.FormatInt(int64(value), 10) }

			m := New[int, string](comparator.OrderedComparator[int])

			for _, value := range xs {
				if value != x {
					m, _ = m.Insert(value, formatter(value))
				}
			}

			newM, affected := m.Delete(x)
			return !affected && newM == m
		},
	})
}
