package client

import (
	"context"
	"math"
	"time"
)

// note that returning either false or a non-nil error will result in the call not being retried
type RetryFunc func(ctx context.Context, req Request, retryCount int, err error) (bool, error)

// RetryAlways always retry on error
func RetryAlways(ctx context.Context, req Request, retryCount int, err error) (bool, error) {
	return true, nil
}

// RetryOnError retries a request on a 500 or timeout error
//func RetryOnError(ctx context.Context, req Request, retryCount int, err error) (bool, error) {
//	if err == nil {
//		return false, nil
//	}
//
//	e := errors.Unwrap(err)
//	if e == nil {
//		return false, nil
//	}
//
//	switch e.Code {
//	// retry on timeout or internal server error
//	case 408, 500:
//		return true, nil
//	default:
//		return false, nil
//	}
//}

// Do is a function x^e multiplied by a factor of 0.1 second.
// Result is limited to 2 minute.
func Do(attempts int) time.Duration {
	if attempts > 13 {
		return 2 * time.Minute
	}
	return time.Duration(math.Pow(float64(attempts), math.E)) * time.Millisecond * 100
}

type BackoffFunc func(ctx context.Context, req Request, attempts int) (time.Duration, error)

func exponentialBackoff(ctx context.Context, req Request, attempts int) (time.Duration, error) {
	return Do(attempts), nil
}
