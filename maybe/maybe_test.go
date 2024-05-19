package maybe

import (
	"reflect"
	"testing"

	immutable_func "github.com/freebirdljj/immutable/func"
	"github.com/freebirdljj/immutable/internal/quick"
)

func TestBind(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
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

func TestMap(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"Map(Nothing(), f) === Nothing()": func() bool {
			return Map(Nothing[int](), immutable_func.Konst[int]("")).IsNothing()
		},
		"Map(Just(x), f) === Just(f(x))": func(x int, y string) bool {
			return reflect.DeepEqual(
				Map(Just(x), immutable_func.Konst[int](y)),
				Just(y),
			)
		},
	})
}

func TestMaybeIsJust(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"Just(x).IsJust() === true": func(x int) bool {
			return Just(x).IsJust()
		},
	})
}

func TestMaybeIsNothing(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
		"Nothing().IsNothing() === true": func() bool {
			return Nothing[int]().IsNothing()
		},
	})
}

func TestMaybeToGoPointer(t *testing.T) {
	quick.CheckProperties(t, map[string]any{
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
	quick.CheckProperties(t, map[string]any{
		"Just(x).OrValue(y) == x": func(x int, y int) bool {
			return Just(x).OrValue(y) == x
		},
		"Nothing().OrValue(x) == x": func(x int) bool {
			return Nothing[int]().OrValue(x) == x
		},
	})
}
