package rateLimitingStrategies

type LimiterStrategy interface {
	addRequest() bool
}

type Limiter struct {
	limiterStrategy LimiterStrategy
}

func NewLimiter(limiterStrategy LimiterStrategy) *Limiter {
	return &Limiter{
		limiterStrategy: limiterStrategy,
	}
}

func (l *Limiter) ChangeLimiterStrategy(limiterStrategy LimiterStrategy) {
	l.limiterStrategy = limiterStrategy
}
func (l *Limiter) AddRequest() bool {
	return l.limiterStrategy.addRequest()
}
