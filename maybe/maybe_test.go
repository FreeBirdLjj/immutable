package maybe

import (
	"reflect"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/require"

	immutable_func "github.com/freebirdljj/immutable/func"
)

func TestBind(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"Bind(Nothing(), f) === Nothing()": func() bool {
			return Bind(Nothing[int](), immutable_func.Konst[int](Just(""))).IsNothing()
		},
		"Bind(Just(x), f) == Nothing() if f(x) == Nothing()": func(x int) bool {
			return Bind(Just(x), immutable_func.Konst[int](Nothing[string]())).IsNothing()
		},
		"Bind(Just(x), f) == f(x) if f(x) returns a `Just`": func(x int, y string) bool {
			return reflect.DeepEqual(
				Bind(Just(x), immutable_func.Konst[int](Just(y))),
				Just(y),
			)
		},
	})
}

func TestMaybeIsJust(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"Just(x).IsJust() === true": func(x int) bool {
			return Just(x).IsJust()
		},
	})
}

func TestMaybeIsNothing(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"Nothing().IsNothing() === true": func() bool {
			return Nothing[int]().IsNothing()
		},
	})
}

func TestMaybeToGoPointer(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"*(Just(x).ToGoPointer()) === x": func(x int) bool {
			return *Just(x).ToGoPointer() == x
		},
		"Nothing().ToGoPointer() === nil": func() bool {
			return Nothing[int]().ToGoPointer() == nil
		},
		"FromGoPointer(ptr).ToGoPointer() === ptr": func(ptr *int) bool {
			return FromGoPointer(ptr).ToGoPointer() == ptr
		},
	})
}

func TestMaybeOrValue(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"Just(x).OrValue(y) == x": func(x int, y int) bool {
			return Just(x).OrValue(y) == x
		},
		"Nothing().OrValue(x) == x": func(x int) bool {
			return Nothing[int]().OrValue(x) == x
		},
	})
}

func checkProperties(t *testing.T, properties map[string]any) {
	for name, property := range properties {
		name, property := name, property
		t.Run(name, func(t *testing.T) {

			t.Parallel()

			err := quick.Check(property, nil)
			require.NoError(t, err)
		})
	}
}
