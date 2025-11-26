# GoRateLimiter

A lightweight, thread-safe rate limiting library for Go. This package provides multiple rate limiting algorithms to control the rate at which requests are processed.

## Features

- 🔒 **Thread-Safe**: Built with synchronization primitives for concurrent access
- ⚡ **Multiple Algorithms**: Fixed Window and No Rate Limiter implementations
- 🎯 **Simple API**: Easy-to-use interface for integrating rate limiting
- 📦 **Lightweight**: Minimal dependencies, production-ready

## Installation

```bash
go get github.com/Jinish9302/GoRateLimiter
```

## Quick Start

### Using FixedWindowRateLimiter

The `FixedWindowRateLimiter` enforces a maximum number of requests within a fixed time window:

```go
package main

import (
	"fmt"
	"time"
	"github.com/Jinish9302/GoRateLimiter/package/ratelimiters"
)

func main() {
	// Allow 10 requests per 1 second
	config := ratelimiters.FixedWindowRateLimiterConfig{
		Limit:          10,
		WindowDuration: time.Second,
	}

	limiter := ratelimiters.NewFixedWindowRateLimiter(config)

	// Check if request is allowed
	allowed, err := limiter.Allow()
	if err != nil {
		panic(err)
	}

	if allowed {
		fmt.Println("Request allowed")
		// Process request
	} else {
		fmt.Println("Rate limit exceeded")
	}

	// Wait until next request is allowed
	err = limiter.Wait()
	if err != nil {
		panic(err)
	}
	fmt.Println("Ready for next request")
}
```

## API Reference

### RateLimiter Interface

All rate limiters implement the `RateLimiter` interface:

```go
type RateLimiter interface {
	Allow() (bool, error)  // Check if request is allowed
	Wait() error           // Wait until next request is allowed
}
```

#### Methods

- **`Allow() (bool, error)`**: Checks if a request is allowed under the current rate limit. Returns `true` if allowed, `false` otherwise.
- **`Wait() error`**: Blocks until the next request is allowed. Useful for controlled request processing.

### FixedWindowRateLimiter

Implements the fixed window rate limiting algorithm. Divides time into fixed intervals (windows) and allows a maximum number of requests per window.

#### Configuration

```go
type FixedWindowRateLimiterConfig struct {
	Limit          int           // Max requests per window
	WindowDuration time.Duration // Length of each window
}
```

#### Usage

```go
config := ratelimiters.FixedWindowRateLimiterConfig{
	Limit:          100,
	WindowDuration: time.Minute,
}

limiter := ratelimiters.NewFixedWindowRateLimiter(config)

// Check if request is allowed
if allowed, err := limiter.Allow(); err == nil && allowed {
	// Process request
}

// Wait for rate limit to reset
limiter.Wait()
```

#### Methods

- **`NewFixedWindowRateLimiter(config FixedWindowRateLimiterConfig)`**: Creates a new fixed window rate limiter
- **`Allow() (bool, error)`**: Checks if request is allowed without blocking
- **`Wait() error`**: Waits until the next request can be made
- **`Reset() error`**: Manually reset the rate limiter

## Examples

### Example 1: HTTP Server with Rate Limiting

```go
package main

import (
	"fmt"
	"net/http"
	"time"
	"github.com/Jinish9302/GoRateLimiter/package/ratelimiters"
)

func main() {
	// 5 requests per second
	limiter := ratelimiters.NewFixedWindowRateLimiter(
		ratelimiters.FixedWindowRateLimiterConfig{
			Limit:          5,
			WindowDuration: time.Second,
		},
	)
	
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		allowed, _ := limiter.Allow()
		if !allowed {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	http.ListenAndServe(":8080", nil)
}
```

### Example 2: Sequential Request Processing

```go
limiter := ratelimiters.NewFixedWindowRateLimiter(
	ratelimiters.FixedWindowRateLimiterConfig{
		Limit:          10,
		WindowDuration: time.Second,
	},
)

for i := 0; i < 100; i++ {
	// Wait until request is allowed
	limiter.Wait()
	
	// Process request
	fmt.Printf("Processing request %d\n", i)
}
```

## Thread Safety

All rate limiters are thread-safe and can be safely used across multiple goroutines. They use `sync.RWMutex` internally to manage concurrent access.

## Testing

Run the test suite:

```bash
go test ./...
```

Run with coverage:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
GoRateLimiter is a lightweight, high-performance Golang library designed to help developers reliably execute tasks while respecting rate limits.
