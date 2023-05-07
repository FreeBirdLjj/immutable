package either_test

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/freebirdljj/immutable/either"
)

func safeAtoi(s string) either.Either[error, int] {
	return either.FromGoResult(strconv.Atoi(s))
}

func safeSqrt(x float64) either.Either[error, float64] {

	if x < 0 {
		err := fmt.Errorf("negative number %f", x)
		return either.Left[float64](err)
	}

	res := math.Sqrt(x)
	return either.Right[error](res)
}

func safeDiv(a int, b int) either.Either[error, int] {

	if b == 0 {
		err := errors.New("divided by zero")
		return either.Left[int](err)
	}

	res := a / b
	return either.Right[error](res)
}

func f(s string) either.Either[error, int] {
	return either.Run(func(computation *either.Computation[error, int]) int {
		var x int = either.Bind(computation, safeAtoi(s))
		var root float64 = either.Bind(computation, safeSqrt(float64(x)))
		var res int = either.Bind(computation, safeDiv(1, int(root)))
		return res
	})
}

func ExampleEither() {
	for _, s := range []string{"abc", "-1", "0", "1"} {
		res, err := either.ToGoResult(f(s))
		if err != nil {
			fmt.Printf("failed to handle `%s`: %s\n", s, err)
			continue
		}
		fmt.Printf("f(`%s`): %d\n", s, res)
	}
	// Output:
	// failed to handle `abc`: strconv.Atoi: parsing "abc": invalid syntax
	// failed to handle `-1`: negative number -1.000000
	// failed to handle `0`: divided by zero
	// f(`1`): 1
}
