package rateLimitingStrategies

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTokenBucket(t *testing.T) {

	testCases := []struct {
		name       string
		bucketSize int
		refillRate int
		requests   int
		expected   []bool
		sleepFunc  func()
	}{
		{
			name:       "OK",
			bucketSize: 10,
			refillRate: 1,
			requests:   10,
			expected:   []bool{true, true, true, true, true, true, true, true, true, true},
		},
		{
			name:       "BucketExceeded",
			bucketSize: 10,
			refillRate: 1,
			requests:   11,
			expected:   []bool{true, true, true, true, true, true, true, true, true, true, false},
		},
		{
			name:       "BucketExceededPeriodically",
			bucketSize: 10,
			refillRate: 1,
			requests:   15,
			expected:   []bool{true, true, true, true, true, true, true, true, true, true, false, true, false, true, false},
			sleepFunc:  func() { time.Sleep(time.Second) },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tb := NewtokenBucket(tc.bucketSize, tc.refillRate)

			for i := 0; i < tc.requests; i++ {
				isAdded := tb.addRequest()
				require.Equal(t, tc.expected[i], isAdded)
				if !isAdded && tc.sleepFunc != nil {
					tc.sleepFunc()
				}
			}
		})
	}

}
