package slipstream

import (
	"context"
	"fmt"
)

type Worker struct {
	broker      Broker
	dispatcher  *Dispatcher
	retryPolicy RetryPolicy
	scheduler   Scheduler
	dlq         DeadLetterQueue
}

func NewWorker(broker Broker, dispatcher *Dispatcher, retryPolicy RetryPolicy, scheduler Scheduler, dlq DeadLetterQueue) *Worker {
	return &Worker{
		broker:      broker,
		dispatcher:  dispatcher,
		retryPolicy: retryPolicy,
		scheduler:   scheduler,
		dlq:         dlq,
	}
}

// Job Acquisition
func (w *Worker) Start(ctx context.Context) {
	for {
		job, err := w.broker.Dequeue(ctx)
		if err != nil {
			return
		}

		if err := w.processJob(ctx, job); err != nil {
			// future logging/metrics
			// Will implement later
		}

	}
}

// Job processing
func (w *Worker) processJob(
	ctx context.Context,
	job *Job,
) error {

	job.Status = JobRunning

	if err := w.dispatcher.Dispatch(ctx, *job); err != nil { //Dispatch the job, if it fails, proceed with retry

		if !w.retryPolicy.ShouldRetry(*job, err) { //Consult the retry policy whether the job should be retried

			job.Status = JobFailed

			if dlqErr := w.dlq.Enqueue(ctx, *job); dlqErr != nil {
				fmt.Println("Could not add a failed job into the dead letter queue")
				return dlqErr
				
			}
		

			return err
		}

		job.RetryAttempt++ //Increment the retry attempt

		job.Status = JobRetrying

		delay := w.retryPolicy.Delay(*job) //Consult the retry policy to determine the delay before retrying

		if err := w.scheduler.Schedule(ctx, *job, delay); err != nil { //Schedule the job to be retried with the determined delay
			return err
		}
		return nil
	}

	if err := w.broker.Ack(ctx, *job); err != nil { //Acknowledge the job to the broker, if it fails, return the error
		return err
	}

	// The job has been successfully acknowledged by the broker.
	// Mark it as completed.
	job.Status = JobCompleted

	return nil
}
