package slipstream

import (
	"context"
	"sync"
)

type WorkerPool struct {
	size        int
	broker      Broker
	dispatcher  *Dispatcher
	retryPolicy RetryPolicy
	scheduler   Scheduler
	dlq         DeadLetterQueue
	wg          sync.WaitGroup
	startOnce   sync.Once
}

func NewWorkerPool(
	size int,
	broker Broker,
	dispatcher *Dispatcher,
	retryPolicy RetryPolicy,
	scheduler Scheduler,
	dlq DeadLetterQueue,
) *WorkerPool {
	return &WorkerPool{
		size:        size,
		broker:      broker,
		dispatcher:  dispatcher,
		retryPolicy: retryPolicy,
		scheduler:   scheduler,
		dlq:         dlq,
	}
}

func (p *WorkerPool) Start(ctx context.Context) {
	p.startOnce.Do(func() {
		p.wg.Add(p.size)

		for i := 0; i < p.size; i++ {
			go func() {
				defer p.wg.Done()

				NewWorker(
					p.broker,
					p.dispatcher,
					p.retryPolicy,
					p.scheduler,
					p.dlq,
				).Start(ctx)
			}()
		}
	})
}

func (p *WorkerPool) Wait(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
