package slipstream

import "time"

//Exponential Backoff: Delay increases exponentially with each retry attempt

type RetryPolicy interface {
	ShouldRetry(job Job, err error) bool //Check whether the job should be retried based on the error
	Delay(job Job) time.Duration //Will use exponential backoff
}

