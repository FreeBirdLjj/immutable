package immutable_map

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/freebirdljj/immutable/comparator"
)

func TestMapInsert(t *testing.T) {

	t.Parallel()

	t.Run("should succeed to insert new value", func(t *testing.T) {

		t.Parallel()

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		kvs := []KeyValuePair[int, string]{
			{Key: 1, Value: "one"},
			{Key: 2, Value: "two"},
			{Key: 3, Value: "three"},
			{Key: 4, Value: "four"},
			{Key: 5, Value: "five"},
			{Key: 6, Value: "six"},
		}

		r.Shuffle(len(kvs), func(i int, j int) {
			kvs[i], kvs[j] = kvs[j], kvs[i]
		})

		m := NewMap[int, string](comparator.OrderedComparator[int])
		for _, kv := range kvs {
			m, _ = m.Insert(kv.Key, kv.Value)
		}

		gotItems := m.KeyValuePairs()

		assert.ElementsMatch(t, kvs, gotItems)
	})
	t.Run("should succeed to update existing value", func(t *testing.T) {

		t.Parallel()

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		kvs := []KeyValuePair[int, string]{
			{Key: 1, Value: "one"},
			{Key: 2, Value: "two"},
			{Key: 3, Value: "three"},
			{Key: 4, Value: "four"},
			{Key: 5, Value: "five"},
			{Key: 6, Value: "six"},
		}
		overriddenKVs := []KeyValuePair[int, string]{
			{Key: 1, Value: "1"},
			{Key: 2, Value: "2"},
			{Key: 3, Value: "3"},
			{Key: 4, Value: "4"},
			{Key: 5, Value: "5"},
			{Key: 6, Value: "6"},
		}

		r.Shuffle(len(kvs), func(i int, j int) {
			kvs[i], kvs[j] = kvs[j], kvs[i]
		})

		m := NewMapFromKeyValuePairs(comparator.OrderedComparator[int], kvs...)

		r.Shuffle(len(overriddenKVs), func(i int, j int) {
			overriddenKVs[i], overriddenKVs[j] = overriddenKVs[j], overriddenKVs[i]
		})

		for _, kv := range overriddenKVs {

			newM, affected := m.Insert(kv.Key, kv.Value)
			assert.False(t, affected)

			gotValue, has := newM.Index(kv.Key)
			assert.True(t, has)
			assert.Equal(t, kv.Value, gotValue)

			assert.Equal(t, len(kvs), newM.Count())
		}
	})
}

func TestMapDelete(t *testing.T) {

	t.Parallel()

	t.Run("should succeed to delete existing key-value pair", func(t *testing.T) {

		t.Parallel()

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		kvs := []KeyValuePair[int, string]{
			{Key: 1, Value: "one"},
			{Key: 2, Value: "two"},
			{Key: 3, Value: "three"},
			{Key: 4, Value: "four"},
			{Key: 5, Value: "five"},
			{Key: 6, Value: "six"},
		}

		r.Shuffle(len(kvs), func(i int, j int) {
			kvs[i], kvs[j] = kvs[j], kvs[i]
		})

		m := NewMapFromKeyValuePairs(comparator.OrderedComparator[int], kvs...)

		for _, kv := range kvs {

			newM, affected := m.Delete(kv.Key)
			assert.True(t, affected)

			_, has := newM.Index(kv.Key)
			assert.False(t, has)

			assert.Equal(t, len(kvs)-1, newM.Count())
		}
	})
	t.Run("should succeed to delete existing key-value pair", func(t *testing.T) {

		t.Parallel()

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		kvs := []KeyValuePair[int, string]{
			{Key: 1, Value: "one"},
			{Key: 2, Value: "two"},
			{Key: 3, Value: "three"},
			{Key: 4, Value: "four"},
			{Key: 5, Value: "five"},
			{Key: 6, Value: "six"},
		}

		r.Shuffle(len(kvs), func(i int, j int) {
			kvs[i], kvs[j] = kvs[j], kvs[i]
		})

		m := NewMapFromKeyValuePairs(comparator.OrderedComparator[int], kvs...)

		for _, kv := range kvs {

			deletedM, _ := m.Delete(kv.Key)

			newM, affected := deletedM.Delete(kv.Key)
			assert.False(t, affected)

			assert.Equal(t, len(kvs)-1, newM.Count())
		}
	})
}
