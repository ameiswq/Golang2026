package retry

import (
	"math/rand"
	"time"
)

func CalculateBackoff(attempt int, baseDelay time.Duration, maxDelay time.Duration) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	backoff := baseDelay
	for i := 0; i < attempt; i++ {
		backoff *= 2
		if backoff > maxDelay {
			backoff = maxDelay
			break
		}
	}
	if backoff <= 0 {
		return 0
	}
	return time.Duration(rand.Int63n(int64(backoff)))
}
