package slipstream

import "context"

type Handler func(ctx context.Context, job Job) error
