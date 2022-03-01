package limiter

import "errors"

type Limiter interface {
	Limit() error
	ChangeQpsThreshold(newQpsThreshold int64)
	Close()
}

const (
	LeakyBucketTyp   = "leak_bucket"
	SildingWindowTyp = "silding_window"
)

func NewLimiter(qpsThreshold int64, limiterTyp string) (Limiter, error) {
	switch limiterTyp {
	case LeakyBucketTyp:
		return NewLeakyBucketRateLimiter(qpsThreshold), nil
	case SildingWindowTyp:
		return NewSlidingWindowRateLimiter(qpsThreshold), nil
	default:
		return nil, errors.New("limiter type error")
	}
}
