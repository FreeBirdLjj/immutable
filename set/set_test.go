package immutable_set

import (
	"testing"

	"github.com/freebirdljj/immutable/comparator"
	"github.com/freebirdljj/immutable/internal/quick"
)

func TestSetInsert(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"should succeed to insert a new value": func(xs []string, x string) bool {

			s := New(comparator.OrderedComparator[string])

			for _, value := range xs {
				if value != x {
					s, _ = s.Insert(value)
				}
			}

			newS, affected := s.Insert(x)
			return affected && newS.Has(x)
		},
		"should succeed to insert an existing value": func(xs []string, x string) bool {

			s := New(comparator.OrderedComparator[string])

			for _, value := range append(xs, x) {
				s, _ = s.Insert(value)
			}

			_, affected := s.Insert(x)
			return !affected
		},
	})
}

func TestSetDelete(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"should succeed to delete an existing value": func(xs []string, x string) bool {

			s := FromValues(comparator.OrderedComparator[string], append(xs, x)...)

			newS, affected := s.Delete(x)
			return affected && !newS.Has(x)
		},
		"should succeed to delete a non-existing value": func(xs []string, x string) bool {

			s := New(comparator.OrderedComparator[string])

			for _, value := range xs {
				if value != x {
					s, _ = s.Insert(value)
				}
			}

			newS, affected := s.Delete(x)
			return !affected && newS == s
		},
	})
}
