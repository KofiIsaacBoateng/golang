package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Concurrent Task Orchestrator with Dynamic Priorities, Timeouts, and Result Consolidation.

type Job struct {
	ID int
	Duration time.Duration
	Priority string // "HIGH" or "LOW"
}

type JobResult struct {
	job Job
	err error
}

// Dispatcher reads raw jobs data and streams them into their respective priority-based channels
func dispatcher(jobs []Job, highPriority chan<- Job, lowPriority chan<- Job) {

	for _, job := range jobs {
		if(job.Priority == "HIGH") {
			highPriority <- job;
		}else {
			lowPriority <- job;
		}
	}


	// Close queues to signal that no job of this priority are coming
	close(highPriority)
	close(lowPriority)
}

func worker(id int, highPriority <-chan Job, lowPriority <-chan Job, resultChannel chan<- JobResult) {
	for {
		var job Job;
		var open bool;

		select {
		case job, open = <- highPriority:
			if(!open){
				highPriority = nil
			}
		default:
			// Do nothing and return to the next line of code
		}


		if(highPriority == nil ){
			job, open = <- lowPriority;
			if(!open){
				break;
			}
		}else if job.ID == 0 {
			// If we didn't pick up a high priority job in the check above, block on both
			select {
			case job, open = <- highPriority:
				if(!open) {
					highPriority = nil;
					continue;
				}
			case job, open = <- lowPriority:
				if(!open) {
					lowPriority = nil;
					continue;
				}
			}
		}


		// Execute the job with a strict SLA timeout of 400ms using context
		ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond);
		
		// process job for the currently received job
		err := processJob(ctx, id, job)
		
		cancel() // Free context resources immediately

		resultChannel <- JobResult{job: job, err: err}

		job = Job{} // reset job allocation tracker for the next iteration
	}
}

// Worker's operation
func processJob(ctx context.Context,workerId int, job Job) error {
	fmt.Printf("[Worker %d] Started Job %d (%s) - Expected: %v\n", workerId, job.ID, job.Priority, job.Duration)

	workDone := make(chan bool, 1);


	go func() {
		time.Sleep(job.Duration); // Simulate processing time
		workDone <- true;
	}()

	select{
	case <- workDone:
		fmt.Printf("\t[Worker %d] Completed Job %d successfully\n", workerId, job.ID)
		return nil
	case <- ctx.Done():
		fmt.Printf("\t[Worker %d] ! TIMEOUT ALERT ! Job %d breached SLA\n", workerId, job.ID)
		return errors.New("timeout: job execution exceeded maximum allowed time SLA")	}
}



func main() {
	// 1. Set up data payload
	jobs := []Job{
		{ID: 1, Duration: 200 * time.Millisecond, Priority: "LOW"},
		{ID: 2, Duration: 600 * time.Millisecond, Priority: "HIGH"}, // Will timeout (> 400ms)
		{ID: 3, Duration: 100 * time.Millisecond, Priority: "HIGH"},
		{ID: 4, Duration: 150 * time.Millisecond, Priority: "LOW"},
		{ID: 5, Duration: 300 * time.Millisecond, Priority: "HIGH"},
	}


	// initiate channels
	highPriorityQueue := make(chan Job, len(jobs));
	lowPriorityQueue := make(chan Job, len(jobs));
	resultChannel := make(chan JobResult, len(jobs));

	// dispatch jobs into their respective priority-based channels
	go dispatcher(jobs, highPriorityQueue, lowPriorityQueue);


	// Launch worker pool (2 concurrent workers )
	noOfWorkers := 2;
	for i := range noOfWorkers {
		go worker(i, highPriorityQueue, lowPriorityQueue, resultChannel);
	}


	// Collect results dynamically
	successCount := 0;
	failureCount := 0;

	for range jobs { // loop automatically until channel is drained!
		result := <-resultChannel;
		if result.err != nil {
			failureCount ++
			continue;
		}

		successCount++
	}

	fmt.Printf("\n--- Orchestration Summary ---\nSuccesses: %d\nFailures/Timeouts: %d\n", successCount, failureCount);
}