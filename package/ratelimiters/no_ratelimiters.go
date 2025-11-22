package ratelimiters

import (
	"time"
)

// NoRateLimiter does not rate limit.
type NoRateLimiter struct{}

// Allow always returns true and nil error, since NoRateLimiter does not rate limit.
func (r *NoRateLimiter) Allow() (bool, error) {
	return true, nil
}

// Wait always returns 0, since NoRateLimiter does not rate limit.
func (r *NoRateLimiter) Wait() (time.Duration, error) {
	return 0, nil
}

// SleepTillAllowed does nothing, since NoRateLimiter does not rate limit.
// It exists solely to match the RateLimiter interface.
func (r *NoRateLimiter) SleepTillAllowed() error {
	time.Sleep(0)
	return nil
}
