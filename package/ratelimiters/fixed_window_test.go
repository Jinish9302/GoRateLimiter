package ratelimiters

import (
	"context"
	"testing"
	"time"
)

func TestFixedWindowRateLimiter_Config(t *testing.T) {
	limiter := NewFixedWindowRateLimiter(FixedWindowRateLimiterConfig{
		Limit:          2,
		WindowDuration: time.Second,
	})
	config := limiter.GetFixedWindowRateLimiterConfig(FixedWindowRateLimiterConfig{})
	if config.Limit != 2 || config.WindowDuration != time.Second {
		t.Errorf("Expected: 2, 1s || Got: %v, %v", config.Limit, config.WindowDuration)
	}
}

func TestFixedWindowRateLimiter_GetFixedWindowRateLimiterConfig(t *testing.T) {
	limiter := NewFixedWindowRateLimiter(FixedWindowRateLimiterConfig{
		Limit:          2,
		WindowDuration: time.Second,
	})
	config := limiter.GetFixedWindowRateLimiterConfig(FixedWindowRateLimiterConfig{})
	if config.Limit != 2 || config.WindowDuration != time.Second {
		t.Errorf("Expected: 2, 1s || Got: %v, %v", config.Limit, config.WindowDuration)
	}
}

func TestFixedWindowRateLimiter_Reset(t *testing.T) {
	limiter := NewFixedWindowRateLimiter(FixedWindowRateLimiterConfig{
		Limit:          2,
		WindowDuration: time.Second,
	})
	allowed, allowErr := limiter.Allow()
	if allowErr != nil || !allowed {
		t.Errorf("Expected: allowed->true,allowErr->nil || Got: allowed->%v,allowErr->%v", allowed, allowErr)
	}
	err := limiter.Reset()
	if err != nil || limiter.requestCount != 0 {
		t.Errorf("Expected: err->nil,limiter.requestCount->0 || Got: err->%v,limiter.requestCount->%v", err, limiter.requestCount)
	}
}

func TestFixedWindowRateLimiter_ResetIfWindowExpired(t *testing.T) {
	limiter := NewFixedWindowRateLimiter(FixedWindowRateLimiterConfig{
		Limit:          2,
		WindowDuration: 100 * time.Millisecond,
	})
	_, _ = limiter.Allow()
	time.Sleep(100 * time.Millisecond)
	err := limiter.ResetIfWindowExpired()
	if err != nil || limiter.requestCount != 0 {
		t.Errorf("Expected: err->nil,limiter.requestCount->0 || Got: err->%v,limiter.requestCount->%v", err, limiter.requestCount)
	}
}

func TestFixedWindowRateLimiter_Allow(t *testing.T) {
	limiter := NewFixedWindowRateLimiter(FixedWindowRateLimiterConfig{
		Limit:          5,
		WindowDuration: time.Second,
	})
	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow()
		if err != nil || !allowed {
			t.Errorf("Expected: true, nil || Got: %v, %v", allowed, err)
		}
	}
	allowed, err := limiter.Allow()
	if err != nil || allowed {
		t.Errorf("Expected: false, nil || Got: %v, %v", allowed, err)
	}
	time.Sleep(time.Second)
	allowed, err = limiter.Allow()
	if err != nil || !allowed {
		t.Errorf("Expected: true, nil || Got: %v, %v", allowed, err)
	}
}

func TestFixedWindowRateLimiter_Wait(t *testing.T) {
	limiter := NewFixedWindowRateLimiter(FixedWindowRateLimiterConfig{
		Limit:          2,
		WindowDuration: 100 * time.Millisecond,
	})
	for i := 0; i < 2; i++ {
		err := limiter.Wait(context.Background())
		if err != nil {
			t.Errorf("Expected: nil || Got: %v", err)
		}
	}
	err := limiter.Wait(context.Background())
	if err != nil {
		t.Errorf("Expected: nil || Got: %v", err)
	}
}

func TestFixedWindowRateLimiter_Remaining(t *testing.T) {
	limiter := NewFixedWindowRateLimiter(FixedWindowRateLimiterConfig{
		Limit:          2,
		WindowDuration: time.Second,
	})
	for i := 0; i <= 2; i++ {
		remaining, _, err := limiter.Remaining()
		expected := 2 - i
		if err != nil || remaining != expected {
			t.Errorf("Expected: err->nil,remaining->%v || Got: err->%v,remaining->%v", expected, err, remaining)
		}
		_, _ = limiter.Allow()
	}
}
