package iter

import (
	"iter"

	"github.com/freebirdljj/immutable/tuple"
)

func SeqFromSeq2[K any, V any](seq2 iter.Seq2[K, V]) iter.Seq[tuple.KeyValuePair[K, V]] {
	return func(yield func(tuple.KeyValuePair[K, V]) bool) {
		seq2(func(key K, value V) bool {
			return yield(tuple.KeyValuePair[K, V]{
				Key:   key,
				Value: value,
			})
		})
	}
}

func Seq2FromSeq[K any, V any](seq iter.Seq[tuple.KeyValuePair[K, V]]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		seq(func(kvPair tuple.KeyValuePair[K, V]) bool {
			return yield(kvPair.Key, kvPair.Value)
		})
	}
}

func Empty[V any]() iter.Seq[V] {
	return func(func(V) bool) {}
}

func Singleton[V any](v V) iter.Seq[V] {
	return func(yield func(V) bool) {
		yield(v)
	}
}

func Repeat[V any](v V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for !yield(v) {
		}
	}
}

// CAUTION: `seq` can't be a single-use iterator nor empty.
func Cycle[V any](seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Map[V1 any, V2 any](seq iter.Seq[V1], f func(V1) V2) iter.Seq[V2] {
	return func(yield func(V2) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}

func Concat[V any](seqs iter.Seq[iter.Seq[V]]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for seq := range seqs {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Take[V any](seq iter.Seq[V], n int) iter.Seq[V] {
	if n == 0 {
		return Empty[V]()
	}
	return func(yield func(V) bool) {
		cnt := 0
		for v := range seq {
			cnt++
			if !yield(v) || cnt == n {
				return
			}
		}
	}
}

func Drop[V any](seq iter.Seq[V], n int) iter.Seq[V] {
	return func(yield func(V) bool) {
		cnt := 0
		for v := range seq {
			cnt++
			if cnt <= n {
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}
