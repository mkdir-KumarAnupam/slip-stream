package slipstream

import "context"

type DeadLetterQueue interface {
	Enqueue(ctx context.Context, job Job) error //Abstracted method for enqueueing jobs to the dead letter queue
}
