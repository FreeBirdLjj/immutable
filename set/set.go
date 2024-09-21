package immutable_set

import (
	"github.com/freebirdljj/immutable/comparator"
	immutable_rb_tree "github.com/freebirdljj/immutable/rb_tree"
)

type (
	// The zero value of `Set` makes nonsense.
	Set[Value any] immutable_rb_tree.RBTree[Value]
)

func New[Value any](cmp comparator.Comparator[Value]) *Set[Value] {
	return (*Set[Value])(immutable_rb_tree.New(cmp))
}

func FromValues[Value any](cmp comparator.Comparator[Value], values ...Value) *Set[Value] {
	return (*Set[Value])(immutable_rb_tree.FromValues(cmp, values...))
}

func (s *Set[Value]) All() func(yield func(value Value) bool) {
	return s.rbTree().All()
}

func (s *Set[Value]) Values() []Value {
	return s.rbTree().Values()
}

func (s *Set[Value]) Empty() bool {
	return s.rbTree().Empty()
}

func (s *Set[Value]) Count() int {
	return s.rbTree().Count()
}

func (s *Set[Value]) Has(value Value) bool {
	return s.rbTree().Lookup(value) != nil
}

// `newSet` returned by `Insert()` is always different from the original one.
// `affected` is true, meaning an actual insertion occurred; otherwise, a replacement occurred.
func (s *Set[Value]) Insert(value Value) (newSet *Set[Value], affected bool) {
	newRBTree, affected := s.rbTree().Insert(value)
	return (*Set[Value])(newRBTree), affected
}

// `affected` is true, meaning that a real deletion occurred, `newSet` will be different from the original;
// otherwise nothing happens, `newSet` is the original one.
func (s *Set[Value]) Delete(value Value) (newSet *Set[Value], affected bool) {
	newRBTree, affected := s.rbTree().Delete(value)
	return (*Set[Value])(newRBTree), affected
}

func (s *Set[Value]) rbTree() *immutable_rb_tree.RBTree[Value] {
	return ((*immutable_rb_tree.RBTree[Value])(s))
}
