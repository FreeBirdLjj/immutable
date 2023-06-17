package either

import (
	"context"
	"runtime"
)

type (
	Either[LeftT any, RightT any] struct {
		left   LeftT
		right  RightT
		isLeft bool
	}

	// Never generate `Computation` value, only use the one passed in as argument `computation` via the callback function `f` of `Run()`.
	Computation[LeftT any, RightT any] struct {
		ch chan Either[LeftT, RightT]
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

// CAUTION: Do not extract from a `Context` that has not had any `Computation` put into it.
func ExtractComputationFromContext[LeftT any, RightT any](ctx context.Context) *Computation[LeftT, RightT] {
	return ctx.Value(computationKey{}).(*Computation[LeftT, RightT])
}

func NewContextWithComputation[LeftT any, RightT any](ctx context.Context, computation *Computation[LeftT, RightT]) context.Context {
	return context.WithValue(ctx, computationKey{}, computation)
}

// CAUTION: Do not invoke `Bind()` with `computation` from another thread.
func Bind[LeftT any, RightT1 any, RightT2 any](computation *Computation[LeftT, RightT1], x Either[LeftT, RightT2]) RightT2 {
	if x.IsLeft() {
		computation.ch <- Left[RightT1](x.Left())
		runtime.Goexit()
	}
	return x.Right()
}

// CAUTION: Do not invoke `BindContext()` with a `Context` that has not had any `Computation` put into it.
func BindContext[RightT1 any, LeftT any, RightT2 any](ctx context.Context, x Either[LeftT, RightT2]) RightT2 {
	computation := ExtractComputationFromContext[LeftT, RightT1](ctx)
	return Bind(computation, x)
}

func Run[LeftT any, RightT any](f func(computation *Computation[LeftT, RightT]) RightT) Either[LeftT, RightT] {

	computation := Computation[LeftT, RightT]{
		ch: make(chan Either[LeftT, RightT], 1),
	}

	go func() {
		right := f(&computation)
		computation.ch <- Right[LeftT](right)
	}()

	res := <-computation.ch
	return res
}

func RunContext[LeftT any, RightT any](ctx context.Context, f func(context.Context) RightT) Either[LeftT, RightT] {
	return Run[LeftT, RightT](func(computation *Computation[LeftT, RightT]) RightT {
		newCtx := NewContextWithComputation(ctx, computation)
		return f(newCtx)
	})
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
