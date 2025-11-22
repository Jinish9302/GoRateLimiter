package ratelimiters

import (
	"time"
)

// RateLimiter defines the interface for rate limiters.
type RateLimiter interface {
	Allow() (bool, error)
	Wait() (time.Duration, error)
	SleepTillAllowed() error
}
