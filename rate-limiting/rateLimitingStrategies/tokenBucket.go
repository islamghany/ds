package rateLimitingStrategies

import (
	"time"
)

// A token bucket has pre-defined capacity(bucketSize). Tokens are put in the bucket at preset rates periodically
// one the bucket is full, no more tokens are added
type tokenBucket struct {
	// maximum number of tokens allowed in the bucket
	bucketSize int
	// number of tokens to add per second
	refillRate int
	lastRefill time.Time
	tokens     int
}

func (tb *tokenBucket) refill() {
	now := time.Now()
	// calculate the number of seconds passed since the last refill
	seconds := now.Sub(tb.lastRefill).Seconds()
	// calculate the number of tokens to add since the last refill
	tokensToAdd := int(seconds) * tb.refillRate
	// add the tokens to the bucket
	tb.tokens = tb.tokens + tokensToAdd
	// ensure the bucket size is not exceeded
	if tb.tokens > tb.bucketSize {
		tb.tokens = tb.bucketSize
	}
	// update the last refill time
	tb.lastRefill = now
}

func (tb *tokenBucket) addRequest() bool {
	tb.refill()
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

func NewtokenBucket(bucketSize, refillRate int) *tokenBucket {
	return &tokenBucket{
		bucketSize: bucketSize,
		refillRate: refillRate,
		lastRefill: time.Now(),
		tokens:     bucketSize,
	}
}
