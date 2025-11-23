package ratelimiters

import (
	"testing"
)

// TestNoRateLimiterAllow tests the correctness of NoRateLimiter's Allow method by checking if it returns the expected value and error.
// It should return true and nil error, since NoRateLimiter does not rate limit.
func TestNoRateLimiterAllow(t *testing.T) {
	rl := NoRateLimiter{}
	allowed, errInAllowed := rl.Allow()
	if errInAllowed != nil || !allowed {
		t.Errorf("Expected: true, nil || Got: %v, %v", allowed, errInAllowed)
	}
}

// TestNoRateLimiterWait tests the correctness of NoRateLimiter's Wait method by checking if it returns the expected value and error.
// It should return 0 and nil error, since NoRateLimiter does not rate limit.
func TestNoRateLimiterWait(t *testing.T) {
	rl := NoRateLimiter{}
	wait, errInWait := rl.Wait()
	if errInWait != nil || wait != 0 {
		t.Errorf("Expected: 0, nil || Got: %v, %v", wait, errInWait)
	}
}

// TestNoRateLimiterSleepTillAllowed tests the correctness of NoRateLimiter's SleepTillAllowed method by checking if it returns the expected value and error.
func TestNoRateLimiterSleepTillAllowed(t *testing.T) {
	rl := NoRateLimiter{}
	errInSleepTillAllowed := rl.SleepTillAllowed()
	if errInSleepTillAllowed != nil {
		t.Errorf("Expected: nil || Got: %v", errInSleepTillAllowed)
	}
}
