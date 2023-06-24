package either_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

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
	return either.Run(func(computation *either.Computation[error]) int {
		var x int = either.Bind(computation, safeAtoi(s))
		var root float64 = either.Bind(computation, safeSqrt(float64(x)))
		var res int = either.Bind(computation, safeDiv(1, int(root)))
		return res
	})
}

func ExampleRun() {
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

func safeHTTPNewRequestWithContext(ctx context.Context, method string, url string, body io.Reader) either.Either[error, *http.Request] {
	return either.FromGoResult(http.NewRequestWithContext(ctx, method, url, body))
}

func safeHTTPClientDo(client *http.Client, req *http.Request) either.Either[error, *http.Response] {
	return either.FromGoResult(client.Do(req))
}

func safeReadAll(reader io.Reader) either.Either[error, []byte] {
	return either.FromGoResult(ioutil.ReadAll(reader))
}

func doSomeIO(ctx context.Context, url string) either.Either[error, string] {
	return either.RunContext[error](ctx, func(ctx context.Context) string {

		client := http.Client{
			Timeout: 10 * time.Second,
		}
		req := either.BindContext(ctx, safeHTTPNewRequestWithContext(ctx, http.MethodGet, url, nil))

		resp := either.BindContext(ctx, safeHTTPClientDo(&client, req))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err := fmt.Errorf("got http response status %d", resp.StatusCode)
			either.BindContext(ctx, either.Left[string](err))
		}

		body := either.BindContext(ctx, safeReadAll(resp.Body))
		return string(body)
	})
}

func ExampleRunContext() {

	failedPath := "/failed"
	successPath := "/success"

	srv := httptest.NewServer(http.HandlerFunc(func(respWriter http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case failedPath:
			respWriter.WriteHeader(http.StatusBadRequest)
		case successPath:
			respWriter.WriteHeader(http.StatusOK)
			io.WriteString(respWriter, "success http response body")
		}
	}))
	defer srv.Close()

	ctx := context.Background()

	for _, path := range []string{
		failedPath,
		successPath,
	} {
		url := srv.URL + path
		res, err := either.ToGoResult(doSomeIO(ctx, url))
		if err != nil {
			fmt.Printf("failed to handle `%s`: %s\n", path, err)
			continue
		}
		fmt.Printf("doSomeIO(`%s`): %s\n", path, res)
	}

	// Output:
	// failed to handle `/failed`: got http response status 400
	// doSomeIO(`/success`): success http response body
}
