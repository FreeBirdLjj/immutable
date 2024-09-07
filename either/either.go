package either

import (
	"context"
	"errors"
	"runtime"

	immutable_func "github.com/freebirdljj/immutable/func"
)

type (
	Either[LeftT any, RightT any] struct {
		left   LeftT
		right  RightT
		isLeft bool
	}

	// Never generate `Computation` value, only use the one passed in as argument `computation` via the callback function `f` of `Run()`.
	Computation[LeftT any] struct {
		returnLeft func(LeftT)
	}

	computationKey struct{}
)

func Left[RightT any, LeftT any](left LeftT) Either[LeftT, RightT] {
	return Either[LeftT, RightT]{
		left:   left,
		isLeft: true,
	}
}

func Right[LeftT any, RightT any](right RightT) Either[LeftT, RightT] {
	return Either[LeftT, RightT]{
		right:  right,
		isLeft: false,
	}
}

func FromGoResult[ResultT any](res ResultT, err error) Either[error, ResultT] {
	if err != nil {
		return Left[ResultT](err)
	}
	return Right[error](res)
}

func ToGoResult[ResultT any](res Either[error, ResultT]) (ResultT, error) {
	return res.Right(), res.Left()
}

func BinaryMap[LeftT any, RightT any, LeftT2 any, RightT2 any](leftMapper func(LeftT) LeftT2, rightMapper func(RightT) RightT2, either Either[LeftT, RightT]) Either[LeftT2, RightT2] {
	if either.IsLeft() {
		return Left[RightT2](leftMapper(either.Left()))
	}
	return Right[LeftT2](rightMapper(either.Right()))
}

func MapLeft[LeftT any, RightT any, LeftT2 any](leftMapper func(LeftT) LeftT2, either Either[LeftT, RightT]) Either[LeftT2, RightT] {
	return BinaryMap(leftMapper, immutable_func.Identity, either)
}

func MapRight[LeftT any, RightT any, RightT2 any](rightMapper func(RightT) RightT2, either Either[LeftT, RightT]) Either[LeftT, RightT2] {
	return BinaryMap(immutable_func.Identity, rightMapper, either)
}

func (either *Either[_, _]) IsLeft() bool {
	return either.isLeft
}

func (either *Either[_, _]) IsRight() bool {
	return !either.isLeft
}

func (either *Either[LeftT, _]) Left() LeftT {
	return either.left
}

func (either *Either[_, RightT]) Right() RightT {
	return either.right
}

func (either *Either[LeftT, RightT]) ToLeft(rightToLeft func(RightT) LeftT) LeftT {
	if either.IsLeft() {
		return either.Left()
	}
	return rightToLeft(either.Right())
}

func (either *Either[LeftT, RightT]) ToRight(leftToRight func(LeftT) RightT) RightT {
	if either.IsRight() {
		return either.Right()
	}
	return leftToRight(either.Left())
}

func (either *Either[LeftT, RightT]) OrLeft(left LeftT) LeftT {
	if either.IsLeft() {
		return either.Left()
	}
	return left
}

func (either *Either[LeftT, RightT]) OrRight(right RightT) RightT {
	if either.IsRight() {
		return either.Right()
	}
	return right
}

// CAUTION: Do not extract from a `Context` that has not had any `Computation` put into it.
func ExtractComputationFromContext[LeftT any](ctx context.Context) *Computation[LeftT] {
	return ctx.Value(computationKey{}).(*Computation[LeftT])
}

func NewContextWithComputation[LeftT any](ctx context.Context, computation *Computation[LeftT]) context.Context {
	return context.WithValue(ctx, computationKey{}, computation)
}

// CAUTION: Do not invoke `Bind()` with `computation` from another thread.
func Bind[LeftT any, RightT any](computation *Computation[LeftT], x Either[LeftT, RightT]) RightT {
	if x.IsLeft() {
		computation.returnLeft(x.Left())
	}
	return x.Right()
}

// CAUTION: Do not invoke `BindContext()` with a `Context` that has not had any `Computation` put into it.
func BindContext[LeftT any, RightT any](ctx context.Context, x Either[LeftT, RightT]) RightT {
	computation := ExtractComputationFromContext[LeftT](ctx)
	return Bind(computation, x)
}

func Run[LeftT any, RightT any](f func(computation *Computation[LeftT]) RightT) Either[LeftT, RightT] {

	ch := make(chan struct{}, 1)
	res := Either[LeftT, RightT]{}

	go func() {
		computation := Computation[LeftT]{
			returnLeft: func(left LeftT) {
				res = Left[RightT](left)
				ch <- struct{}{}
				runtime.Goexit()
			},
		}
		right := f(&computation)
		res = Right[LeftT](right)
		ch <- struct{}{}
	}()

	<-ch
	return res
}

func RunContext[LeftT any, RightT any](ctx context.Context, f func(context.Context) RightT) Either[LeftT, RightT] {
	return Run[LeftT, RightT](func(computation *Computation[LeftT]) RightT {
		newCtx := NewContextWithComputation(ctx, computation)
		return f(newCtx)
	})
}

// Inject a `Computation` into `ctx` if it doesn't have one (but doesn't check if the specific type parameters match).
// Otherwise invoke `f()` directly with the given `ctx`.
func RunPossibleContext[LeftT any, RightT any](ctx context.Context, f func(context.Context) RightT) Either[LeftT, RightT] {
	if ctx.Value(computationKey{}) == nil {
		return RunContext[LeftT](ctx, f)
	}
	right := f(ctx)
	return Right[LeftT](right)
}

func PartitionEithers[LeftT any, RightT any](xs ...Either[LeftT, RightT]) ([]LeftT, []RightT) {
	lefts := make([]LeftT, 0, len(xs)/2)
	rights := make([]RightT, 0, len(xs)/2)
	for _, x := range xs {
		if x.IsLeft() {
			lefts = append(lefts, x.Left())
		} else {
			rights = append(rights, x.Right())
		}
	}
	return lefts, rights
}

func JoinResults[RightT any](xs ...Either[error, RightT]) Either[error, []RightT] {
	lefts, rights := PartitionEithers(xs...)
	if len(lefts) == 0 {
		return Right[error](rights)
	}
	return Left[[]RightT](errors.Join(lefts...))
}
