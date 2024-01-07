package maybe

import (
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/require"
)

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
