package main

import (
	"fmt"
	"islamghany/ds/rate-limiting/rateLimitingStrategies"
	"time"
)

func main() {
	tokenBucketStrategy := rateLimitingStrategies.NewtokenBucket(10, 1)
	limiter := rateLimitingStrategies.NewLimiter(tokenBucketStrategy)

	// test limiter with 10 requests
	for i := 0; i < 10; i++ {
		fmt.Println(limiter.AddRequest())
	}
	time.Sleep(3 * time.Second)
	// test limiter with 10 requests
	for i := 0; i < 10; i++ {
		isBucketExceeded := limiter.AddRequest()
		fmt.Println(isBucketExceeded)
		if !isBucketExceeded {
			time.Sleep(3 * time.Second)
		}

	}

}
