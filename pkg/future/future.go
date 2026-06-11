package future

import (
	"fmt"
	"time"

	"github.com/Meduzz/helper/fp"
)

var (
	ErrTimeout = fmt.Errorf("timeout")
)

// Dont look, nothing to see here.
func Future[T, K any](op fp.Operation[T, K]) func(T) (K, error) {
	timeout := time.NewTimer(250 * time.Millisecond)
	success := make(chan K, 1)
	failure := make(chan error, 1)

	return func(data T) (K, error) {
		go caller(op, data, success, failure)

		defer timeout.Stop()
		defer close(success)
		defer close(failure)

		select {
		case <-timeout.C:
			return *new(K), ErrTimeout
		case it := <-success:
			return it, nil
		case err := <-failure:
			return *new(K), err
		}
	}
}

func caller[T, K any](op fp.Operation[T, K], data T, s chan K, e chan error) {
	value, err := op(data)

	if err != nil {
		e <- err
	} else {
		s <- value
	}
}
