package slipstream

import "context"

type MemoryBroker struct {
	jobs chan Job
}

func NewMemoryBroker(bufferSize int) *MemoryBroker {
	return &MemoryBroker{
		jobs: make(chan Job, bufferSize),
	}
}

func (b *MemoryBroker) Enqueue(ctx context.Context, job Job) error {
	select {
	case <-ctx.Done():
		return ctx.Err() //Returning the context error, which can be used to cancel the operation
	case b.jobs <- job: //Receiving the job on the channel
		return nil //Returning nil to indicate success
	}
	
}

func (b *MemoryBroker) Dequeue(ctx context.Context) (*Job, error) {
	select {
	case <-ctx.Done(): //Returning the context error, which can be used to cancel the operation
		return nil, ctx.Err()
	case job := <-b.jobs: //Returning the job, which can be used to process it
		return &job, nil
	}
}

func (b *MemoryBroker) Ack(ctx context.Context, job Job) error {
	// No-op for the in-memory broker.
	// Once a job is received from the channel,
	// it is considered acknowledged.
	return nil
}
