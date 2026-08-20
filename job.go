	package slipstream
	
	import "github.com/google/uuid"
	
	type JobStatus string
	
	const (
		JobQueued    JobStatus = "queued"
		JobRunning   JobStatus = "running"
		JobCompleted JobStatus = "completed"
		JobFailed    JobStatus = "failed"
		JobRetrying  JobStatus = "retrying"
	)
	
	type Job struct {
		ID           uuid.UUID
		Type         string
		Payload      []byte //Payload will be serialized as bytes from JSON, and then deserialized upon consummption
		Status       JobStatus
		RetryAttempt int
	}
