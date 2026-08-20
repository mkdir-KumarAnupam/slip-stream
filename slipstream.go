package slipstream

import (
	"context"
	"sync"
	"time"
)

type Slipstream struct {
	broker     Broker
	registry   *registry
	workerPool *WorkerPool
	cancel     context.CancelFunc
	syncOnce   sync.Once
}

func New(
	broker Broker,
	workerCount int,
) *Slipstream {

	if broker == nil {
		panic("slipstream: broker cannot be nil")
	}

	if workerCount <= 0 {
		panic("slipstream: workerCount must be greater than 0")
	}

	registry := newRegistry()

	dispatcher := NewDispatcher(registry)

	retryPolicy := NewExponentialRetryPolicy(
		3,
		time.Second,
	)

	scheduler := NewMemoryScheduler(broker)

	dlq := NewMemoryDeadLetterQueue(100)

	workerPool := NewWorkerPool(
		workerCount,
		broker,
		dispatcher,
		retryPolicy,
		scheduler,
		dlq,
	)

	return &Slipstream{
		broker:     broker,
		registry:   registry,
		workerPool: workerPool,
	}
}

func (s *Slipstream) Register(
	jobType string,
	handler Handler,
) error {
	return s.registry.Register(jobType, handler)
}

func (s *Slipstream) Enqueue(
	ctx context.Context,
	job Job,
) error {
	return s.broker.Enqueue(ctx, job)
}

func (s *Slipstream) Start(ctx context.Context) {
	s.syncOnce.Do(func() { //Making sure the start only happens once
		workerCtx, cancel := context.WithCancel(ctx)
		s.cancel = cancel
		s.workerPool.Start(workerCtx)
	})
} //Serializable

func (s *Slipstream) Wait(ctx context.Context) error {
	return s.workerPool.Wait(ctx)
}

func (s *Slipstream) Shutdown(ctx context.Context) error {
	if s.cancel == nil {
		return nil
	}

	s.cancel()
	return s.workerPool.Wait(ctx)
}
