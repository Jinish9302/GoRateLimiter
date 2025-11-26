package ratelimiters

// RateLimiter defines the interface for rate limiters.
type RateLimiter interface {
	Allow() (bool, error)
	Wait() error
}
