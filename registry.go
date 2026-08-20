package slipstream

type registry struct {
	handlers map[string]Handler
}

func newRegistry() *registry {
	return &registry{
		handlers: make(map[string]Handler),
	}
}

func (r *registry) Register(jobType string, handler Handler) error {

	//Validators
	if jobType == "" {
		return ErrInvalidJobType
	}

	if handler == nil {
		return ErrNilHandler
	}

	if _, exists := r.handlers[jobType]; exists {
		return ErrHandlerAlreadyRegistered
	}

	//Register the handler into the map registry
	r.handlers[jobType] = handler

	return nil
}

func (r *registry) Get(jobType string) (Handler, bool) {
	//Simple map fetch
	handler, ok := r.handlers[jobType]
	return handler, ok
}
