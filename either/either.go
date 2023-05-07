package either

import (
	"runtime"
)

type (
	Either[LeftT any, RightT any] struct {
		left   LeftT
		right  RightT
		isLeft bool
	}
	Computation[LeftT any, RightT any] struct {
		ch chan Either[LeftT, RightT]
	}
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

func Bind[LeftT any, RightT1 any, RightT2 any](computation *Computation[LeftT, RightT1], x Either[LeftT, RightT2]) RightT2 {
	if x.IsLeft() {
		computation.ch <- Left[RightT1](x.Left())
		runtime.Goexit()
	}
	return x.Right()
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
