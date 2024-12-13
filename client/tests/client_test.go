package test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/kevicai/job-status-api/client"
)

func TestClient_GetJobStatus(t *testing.T) {
	numJobs := 10

	fmt.Printf("Testing the creation and completion of %v jobs ...\n", numJobs)

	startTimes := sync.Map{}

	var wg sync.WaitGroup
	wg.Add(numJobs) // add 10 jobs to wait for completion

	// Create a JobClient instance
	c := client.NewJobClient()

	handleComplete := func(job *client.Job) {
		defer wg.Done() // mark the task as done

		timeStarted, ok := startTimes.Load(job.ID)
		if !ok {
			t.Log("Failed to get start time for job", job.ID)
		}

		fmt.Println("Job completed:", job.ID, "| Status:", job.Status, "| Time taken:", time.Since(timeStarted.(time.Time)))
	}

	// create jobs
	for i := 0; i < numJobs; i++ {
		jobID, err := c.CreateJob(handleComplete)
		if err != nil {
			t.Error("Failed to create job:", err)
			return
		}

		if jobID == "" {
			t.Error("Expected job ID, got empty string")
			return
		}
		startTimes.Store(jobID, time.Now())
		fmt.Println("Created job with ID:", jobID)
	}

	// wait for all jobs to complete
	wg.Wait()
}

func TestClient_GetJobStatusAfterTrain(t *testing.T) {
	numJobs := 10

	fmt.Printf("Testing the creation and completion of %v jobs after training ...\n", numJobs)

	startTimes := sync.Map{}

	var wg sync.WaitGroup
	wg.Add(numJobs) // add 10 jobs to wait for completion

	// Create a JobClient instance
	c := &client.JobClient{}

	handleComplete := func(job *client.Job) {
		defer wg.Done() // mark the task as done

		timeStarted, ok := startTimes.Load(job.ID)
		if !ok {
			t.Log("Failed to get start time for job", job.ID)
		}

		fmt.Println("Job completed:", job.ID, "| Status:", job.Status, "| Time taken:", time.Since(timeStarted.(time.Time)))
	}

	// create jobs
	for i := 0; i < numJobs; i++ {
		jobID, err := c.CreateJob(handleComplete)
		if err != nil {
			t.Error("Failed to create job:", err)
			return
		}

		if jobID == "" {
			t.Error("Expected job ID, got empty string")
			return
		}
		startTimes.Store(jobID, time.Now())
		fmt.Println("Created job with ID:", jobID)
	}

	// wait for all jobs to complete
	wg.Wait()
}
