package slipstream

import (
	"context"
	"time"
)

type MemoryScheduler struct {
	broker Broker
}

func NewMemoryScheduler(
	broker Broker,
) *MemoryScheduler {
	return &MemoryScheduler{
		broker: broker,
	}
}

func (s *MemoryScheduler) Schedule(
	ctx context.Context,
	job Job,
	delay time.Duration,
) error {

	go func() {
		timer := time.NewTimer(delay)
		defer timer.Stop()

		select {
		case <-timer.C:
			_ = s.broker.Enqueue(ctx, job)

		case <-ctx.Done():
			return
		}
	}()

	return nil
}
