package immutable_map

import (
	"github.com/freebirdljj/immutable/comparator"
	immutable_func "github.com/freebirdljj/immutable/func"
	immutable_rb_tree "github.com/freebirdljj/immutable/rb_tree"
)

type (
	KeyValuePair[Key any, Value any] struct {
		Key   Key
		Value Value
	}

	Map[Key any, Value any] immutable_rb_tree.RBTree[KeyValuePair[Key, Value]]
)

func New[Key any, Value any](cmp comparator.Comparator[Key]) *Map[Key, Value] {
	return (*Map[Key, Value])(immutable_rb_tree.New(
		func(l KeyValuePair[Key, Value], r KeyValuePair[Key, Value]) int {
			return cmp(l.Key, r.Key)
		},
	))
}

func FromGoMap[Key comparable, Value any](cmp comparator.Comparator[Key], goMap map[Key]Value) *Map[Key, Value] {
	m := New[Key, Value](cmp)
	for k, v := range goMap {
		m, _ = m.Insert(k, v)
	}
	return m
}

func FromKeyValuePairs[Key any, Value any](cmp comparator.Comparator[Key], kvPairs ...KeyValuePair[Key, Value]) *Map[Key, Value] {
	return (*Map[Key, Value])(immutable_rb_tree.FromValues(
		func(l KeyValuePair[Key, Value], r KeyValuePair[Key, Value]) int {
			return cmp(l.Key, r.Key)
		},
		kvPairs...,
	))
}

func ToGoMap[Key comparable, Value any](m *Map[Key, Value]) map[Key]Value {
	goMap := make(map[Key]Value, m.Count())
	for _, kvPair := range m.KeyValuePairs() {
		goMap[kvPair.Key] = kvPair.Value
	}
	return goMap
}

func (m *Map[Key, Value]) Empty() bool {
	return m.rbTree().Empty()
}

func (m *Map[Key, Value]) Count() int {
	return m.rbTree().Count()
}

func (m *Map[Key, Value]) Index(key Key) (value Value, has bool) {
	kv := m.rbTree().Lookup(KeyValuePair[Key, Value]{
		Key: key,
	})
	if kv == nil {
		return immutable_func.Zero[Value](), false
	}
	return kv.Value, true
}

func (m *Map[Key, Value]) All() func(yield func(key Key, value Value) bool) {
	rbTreeIter := m.rbTree().All()
	return func(yield func(key Key, value Value) bool) {
		rbTreeIter(func(kvPair KeyValuePair[Key, Value]) bool {
			return yield(kvPair.Key, kvPair.Value)
		})
	}
}

func (m *Map[Key, Value]) KeyValuePairs() []KeyValuePair[Key, Value] {
	return m.rbTree().Values()
}

func (m *Map[Key, Value]) Insert(key Key, value Value) (newMap *Map[Key, Value], affected bool) {
	newRBTree, affected := m.rbTree().Insert(KeyValuePair[Key, Value]{
		Key:   key,
		Value: value,
	})
	return (*Map[Key, Value])(newRBTree), affected
}

func (m *Map[Key, Value]) Delete(key Key) (newMap *Map[Key, Value], affected bool) {
	newRBTree, affected := m.rbTree().Delete(KeyValuePair[Key, Value]{
		Key: key,
	})
	return (*Map[Key, Value])(newRBTree), affected
}

func (m *Map[Key, Value]) rbTree() *immutable_rb_tree.RBTree[KeyValuePair[Key, Value]] {
	return (*immutable_rb_tree.RBTree[KeyValuePair[Key, Value]])(m)
}
