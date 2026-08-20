package slipstream

import (
	"context"
	"time"
)

type Scheduler interface {
	Schedule(
		ctx context.Context,
		job Job,
		delay time.Duration,
	) error
}
