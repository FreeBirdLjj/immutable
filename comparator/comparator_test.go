package comparator

import (
	"testing"

	immutable_func "github.com/freebirdljj/immutable/func"
	"github.com/freebirdljj/immutable/internal/quick"
)

func TestCascadeComparator(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"CascadeComparator(cmp, id) === cmp": func(l int, r int) bool {
			cmp := OrderedComparator[int]
			id := immutable_func.Identity[int]
			return CascadeComparator(cmp, id)(l, r) == cmp(l, r)
		},
		"CascadeComparator(cmp, neg) should return opposite of `cmp` result": func(l int, r int) bool {
			cmp := OrderedComparator[int]
			neg := func(x int) int { return -x }
			return CascadeComparator(cmp, neg)(l, r) == -OrderedComparator(l, r)
		},
	})
}
