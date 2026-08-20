package slipstream

import "context"

type MemoryDeadLetterQueue struct {
	jobs chan Job
}

func NewMemoryDeadLetterQueue(bufferSize int) *MemoryDeadLetterQueue {
	return &MemoryDeadLetterQueue{
		jobs: make(chan Job, bufferSize),
	}
}

func (q *MemoryDeadLetterQueue) Enqueue(ctx context.Context, job Job) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case q.jobs <- job:
		return nil
	}
}


