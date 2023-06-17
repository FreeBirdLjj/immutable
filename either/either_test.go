package either

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/require"
)

func TestComputationInterconvertion(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"computation -> ctx -> computation": func() bool {
			ctx := context.Background()
			computation := new(Computation[error, int])
			newCtx := NewContextWithComputation(ctx, computation)
			gotComputation := ExtractComputationFromContext[error, int](newCtx)
			return gotComputation == computation
		},
	})
}

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

func TestRunContext(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"`RunContext()` returns left if any `BindContext()` receives a left": func(left int, right string) bool {
			res := RunContext[int](context.Background(), func(ctx context.Context) string {
				x := BindContext[string](ctx, Left[string](left))
				return x
			})
			return res.IsLeft() && res.Left() == left
		},
		"`RunContext()` returns left if all `BindContext()`s receive rights": func(left int, right string) bool {
			res := RunContext[int](context.Background(), func(ctx context.Context) string {
				x := BindContext[string](ctx, Right[int](right))
				return x
			})
			return res.IsRight() && res.Right() == right
		},
	})
}

func TestPartitionEithers(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"PartitionEithers(lefts.map(Left) ++ rights.map(Right)) == (lefts, rights)": func(lefts []int, rights []string) bool {

			eithers := make([]Either[int, string], 0, len(lefts)+len(rights))
			for _, left := range lefts {
				eithers = append(eithers, Left[string](left))
			}
			for _, right := range rights {
				eithers = append(eithers, Right[int](right))
			}

			gotLefts, gotRights := PartitionEithers(eithers...)
			return slicesEqual(gotLefts, lefts) && slicesEqual(gotRights, rights)
		},
	})
}

func TestJoinResults(t *testing.T) {

	t.Parallel()

	checkProperties(t, map[string]any{
		"JoinResults(rights.map(Right)) == Right(rights)": func(rights []int) bool {
			results := make([]Either[error, int], len(rights))
			for i, right := range rights {
				results[i] = Right[error](right)
			}
			gotAccumulation := JoinResults(results...)
			return gotAccumulation.IsRight() && slicesEqual(gotAccumulation.Right(), rights)
		},
		"JoinResults(rights.map(Right).append(err.map(Left))) is Left": func(rights []int) bool {

			err := errors.New("error1")

			results := make([]Either[error, int], len(rights)+1)
			for i, right := range rights {
				results[i] = Right[error](right)
			}
			results[len(rights)] = Left[int](err)

			gotAccumulation := JoinResults(results...)
			return gotAccumulation.IsLeft() && errors.Is(gotAccumulation.Left(), err)
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

func slicesEqual[T any](v1 []T, v2 []T) bool {
	return (len(v1) == 0 && len(v2) == 0) || reflect.DeepEqual(v1, v2)
}
