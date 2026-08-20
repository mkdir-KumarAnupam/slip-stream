package slipstream

import "errors"

var (
	ErrInvalidJobType           = errors.New("slipstream: invalid job type")
	ErrNilHandler               = errors.New("slipstream: handler cannot be nil")
	ErrHandlerAlreadyRegistered = errors.New("slipstream: handler already registered for job type")
	ErrHandlerNotFound          = errors.New("slipstream: handler not found")
)
