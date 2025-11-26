package ratelimiters

import (
	"context"
	"sync"
	"time"
)

type FixedWindowRateLimiterConfig struct {
	Limit          int
	WindowDuration time.Duration
}

type FixedWindowRateLimiter struct {
	mux                sync.RWMutex
	limit              int
	windowDuration     time.Duration
	currentWindowStart time.Time
	requestCount       int
}

func NewFixedWindowRateLimiter(config FixedWindowRateLimiterConfig) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		limit:              config.Limit,
		windowDuration:     config.WindowDuration,
		currentWindowStart: time.Now(),
		requestCount:       0,
	}
}

func (limiter *FixedWindowRateLimiter) GetFixedWindowRateLimiterConfig() FixedWindowRateLimiterConfig {
	limiter.mux.RLock()
	defer limiter.mux.RUnlock()
	return FixedWindowRateLimiterConfig{Limit: limiter.limit, WindowDuration: limiter.windowDuration}
}

func (limiter *FixedWindowRateLimiter) unsafeReset() error {
	limiter.currentWindowStart = time.Now()
	limiter.requestCount = 0
	return nil
}
func (limiter *FixedWindowRateLimiter) Reset() error {
	limiter.mux.Lock()
	defer limiter.mux.Unlock()
	err := limiter.unsafeReset()
	return err
}

func (limiter *FixedWindowRateLimiter) unsafeResetIfWindowExpired() error {
	if time.Since(limiter.currentWindowStart) >= limiter.windowDuration {
		return limiter.unsafeReset()
	}
	return nil
}
func (limiter *FixedWindowRateLimiter) ResetIfWindowExpired() error {
	limiter.mux.Lock()
	defer limiter.mux.Unlock()
	err := limiter.unsafeResetIfWindowExpired()
	return err
}

func (limiter *FixedWindowRateLimiter) Allow() (bool, error) {
	limiter.mux.Lock()
	defer limiter.mux.Unlock()
	err := limiter.unsafeResetIfWindowExpired()
	if err == nil && limiter.requestCount < limiter.limit {
		limiter.requestCount++
		return true, nil
	}
	return false, err
}

func (limiter *FixedWindowRateLimiter) Wait(ctx context.Context) error {
	if allowed, _ := limiter.Allow(); allowed {
		return nil
	}
	limiter.mux.RLock()
	timeToWait := time.Until(limiter.currentWindowStart.Add(limiter.windowDuration))
	limiter.mux.RUnlock()
	timer := time.NewTimer(timeToWait)
	defer timer.Stop()
	select {
	case <-timer.C:
		if allowed, _ := limiter.Allow(); allowed {
			return nil
		}
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (limiter *FixedWindowRateLimiter) Remaining() (int, time.Duration, error) {
	limiter.mux.RLock()
	defer limiter.mux.RUnlock()
	err := limiter.unsafeResetIfWindowExpired()
	if err != nil {
		return 0, 0, err
	}
	remaining := limiter.limit - limiter.requestCount
	timeToReset := time.Until(limiter.currentWindowStart.Add(limiter.windowDuration))
	return remaining, timeToReset, nil
}
