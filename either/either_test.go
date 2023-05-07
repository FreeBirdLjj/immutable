package either

import (
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"`run()` returns left if any `bind()` receives a left": func(left int, right string) bool {
			res := Run(func(computation *Computation[int, string]) string {
				x := Bind(computation, Left[string](left))
				return x
			})
			return res.IsLeft() && res.Left() == left
		},
		"`run()` returns left if all `bind()`s receive rights": func(left int, right string) bool {
			res := Run(func(computation *Computation[int, string]) string {
				x := Bind(computation, Right[int](right))
				return x
			})
			return res.IsRight() && res.Right() == right
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
