package ratelimiters

import (
	"context"
	"testing"
	"time"
)

func TestFixedWindowRateLimiter_Config(t *testing.T) {
	testSuite := []struct {
		limit          int
		windowDuration time.Duration
	}{
		{2, time.Second},
		{3, time.Minute},
		{5, time.Hour},
	}
	for i, test := range testSuite {
		limiter := NewFixedWindowRateLimiter(FixedWindowRateLimiterConfig{
			Limit:          test.limit,
			WindowDuration: test.windowDuration,
		})
		if limiter.limit != test.limit || limiter.windowDuration != test.windowDuration {
			t.Errorf("Test Case: %v || Expected: %v, %v || Got: %v, %v", i, test.limit, test.windowDuration, limiter.limit, limiter.windowDuration)
		}
	}
}

func TestFixedWindowRateLimiter_GetFixedWindowRateLimiterConfig(t *testing.T) {
	testSuite := []struct {
		limit          int
		windowDuration time.Duration
	}{
		{2, time.Second},
		{3, time.Minute},
		{5, time.Hour},
	}
	for i, test := range testSuite {
		limiter := NewFixedWindowRateLimiter(FixedWindowRateLimiterConfig{
			Limit:          test.limit,
			WindowDuration: test.windowDuration,
		})
		if limiter.limit != test.limit || limiter.windowDuration != test.windowDuration {
			t.Errorf(
				"Test Case: %v || Expected:%v, %v || Got: %v, %v",
				i,
				test.limit,
				test.windowDuration,
				limiter.limit,
				limiter.windowDuration,
			)
		}
	}
}

func TestFixedWindowRateLimiter_Reset(t *testing.T) {
	limiter := NewFixedWindowRateLimiter(FixedWindowRateLimiterConfig{
		Limit:          2,
		WindowDuration: 100 * time.Millisecond,
	})
	initialTime := limiter.currentWindowStart
	allowed, allowErr := limiter.Allow()
	if (allowErr != nil || !allowed || limiter.currentWindowStart != initialTime) {
		t.Errorf(
			"Test Case: 0 || Expected: allowed->true,allowErr->nil,limiter.currentWindowStart->%v || Got: allowed->%v,allowErr->%v,limiter.currentWindowStart->%v",
			initialTime,
			allowed,
			allowErr,
			limiter.currentWindowStart,
		)
	}
	err := limiter.Reset()
	if (err != nil || limiter.requestCount != 0 || limiter.currentWindowStart!= initialTime) {
		t.Errorf(
			"Test Case: 1 || Expected: err->nil,limiter.requestCount->0,limiter.currentWindowStart->(not %v) || Got: err->%v,limiter.requestCount->%v, limiter.currentWindowStart->%v",
			initialTime,
			err,
			limiter.requestCount,
			limiter.currentWindowStart,
		)
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
		t.Errorf(
			"Expected: err->nil,limiter.requestCount->0 || Got: err->%v,limiter.requestCount->%v",
			err,
			limiter.requestCount,
		)
	}
}

func TestFixedWindowRateLimiter_Allow(t *testing.T) {
	limiter := NewFixedWindowRateLimiter(FixedWindowRateLimiterConfig{
		Limit:          2,
		WindowDuration: 100 * time.Millisecond,
	})
	testSuite := []struct {
		allowed                    bool
		err                        error
		sleepDurationBeforeRequest time.Duration
	}{
		{true, nil, 0},
		{true, nil, 0},
		{false, nil, 0},
		{true, nil, 100 * time.Millisecond},
	}
	for i, test := range testSuite {
		time.Sleep(test.sleepDurationBeforeRequest)
		if allowed, err := limiter.Allow(); err != test.err || allowed != test.allowed {
			t.Errorf("Test Case: %v || Expected: %v, %v || Got: %v, %v", i, test.allowed, test.err, allowed, err)
		}
	}
}

func TestFixedWindowRateLimiter_Wait(t *testing.T) {
	limiter := NewFixedWindowRateLimiter(FixedWindowRateLimiterConfig{
		Limit:          1,
		WindowDuration: 100 * time.Millisecond,
	})
	testSuitse := []struct {
		err error
		interupted bool
	} {
		{nil, false},
		{nil, false},
		{context.Canceled, true},
	}
	for i, test := range testSuitse {
		ctx, cancel := context.WithCancel(context.Background())
		if test.interupted {
			go func() {
				time.Sleep(50 * time.Millisecond)
				cancel()
			}()
		} else {
			defer cancel()
		}
		err := limiter.Wait(ctx)
		if err != test.err {
			t.Errorf(
				"Test Case: %v || Expected: %v || Got: %v",
				i,
				test.err,
				err,
			)
		}
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
