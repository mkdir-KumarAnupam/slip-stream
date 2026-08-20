package slipstream

import "time"

type ExponentialRetryPolicy struct {
	MaxRetries int
	BaseDelay  time.Duration
}

func NewExponentialRetryPolicy(
	maxRetries int,
	baseDelay time.Duration,
) *ExponentialRetryPolicy {

	return &ExponentialRetryPolicy{
		MaxRetries: maxRetries,
		BaseDelay:  baseDelay,
	}
}

func (p *ExponentialRetryPolicy) ShouldRetry(
	job Job,
	err error,
) bool {
	return job.RetryAttempt < p.MaxRetries
}

func (p *ExponentialRetryPolicy) Delay(
	job Job,
) time.Duration {

	// 1st retry -> BaseDelay
	// 2nd retry -> BaseDelay * 2
	// 3rd retry -> BaseDelay * 4
	// ...
	

	//Implementation of exponential backoff using bit shifting
	return p.BaseDelay << job.RetryAttempt //Pretty much: BaseDelay × 2^RetryAttempt
	
}
