package slipstream

import "context"

type Dispatcher struct {
	registry *registry
}

func NewDispatcher(registry *registry) *Dispatcher {
	return &Dispatcher{
		registry: registry,
	}
}

func (d *Dispatcher) Dispatch(ctx context.Context, job Job) error {
	handler, ok := d.registry.Get(job.Type)
	if !ok {
		return ErrHandlerNotFound
	}
	return handler(ctx, job)
}
