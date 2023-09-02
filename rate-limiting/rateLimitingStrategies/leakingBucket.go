package rateLimitingStrategies

// the Leaking Bucket is similar to the token bucket except that requests are processed at fixed rate.
// it is implemented with a queue.
// So when request arrive we check if the queue is full, if it is not full, the request is added to the queue.
// OW the request dropped, the the requests are pulled from the queue and processed at regular interval.
type leakingBucket struct {
	// maximum number of tokens allowed in the bucket
	bucketSize int
	// how many requests can be processed at a fixed rate in seconds.
	outflowRate int
}
