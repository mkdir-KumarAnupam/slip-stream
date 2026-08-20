package slipstream

import "context"

type Broker interface {
	Enqueue(ctx context.Context, job Job) error
	Dequeue(ctx context.Context) (*Job, error)
	Ack(ctx context.Context, job Job) error
}
