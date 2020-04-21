package ratelimit

import (
	"google.golang.org/grpc"
	"myth/go-essential/net/rpc/warden"
)

// alwaysPassLimiter is an example limiter which implements Limiter interface.
// It does not limit any request because Limit function always returns false.
type alwaysPassLimiter struct{}

func (*alwaysPassLimiter) Limit() bool {
	return false
}

// Simple example of server initialization code.
func Example() {
	// Create unary/stream rateLimiters, based on token bucket here.
	// You can implement your own ratelimiter for the interface.
	limiter := &alwaysPassLimiter{}
	_ = grpc.NewServer(
		warden.WithUnaryServerChain(
			UnaryServerInterceptor(limiter),
		),
		warden.WithStreamServerChain(
			StreamServerInterceptor(limiter),
		),
	)
}
