package test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/kevicai/job-status-api/translation"
)

func TestClient_GetJobStatus(t *testing.T) {
	fmt.Println("Testing the creation and completion of 3 jobs ...")

	// Create a JobClient instance
	client := &translation.JobClient{}

	numJobs := 3

	var wg sync.WaitGroup
	wg.Add(numJobs) // add 3 jobs to wait for completion

	handleComplete := func(job *translation.Job) {
		defer wg.Done() // mark the task as done
		fmt.Println("Job completed:", job.ID, "Status:", job.Status)
	}

	// create 3 jobs
	for i := 0; i < numJobs; i++ {
		jobID, err := client.CreateJob(handleComplete)
		if err != nil {
			t.Error("Failed to create job:", err)
			return
		}

		if jobID == "" {
			t.Error("Expected job ID, got empty string")
			return
		}

		fmt.Println("Created job with ID:", jobID)
	}

	// wait for all jobs to complete
	wg.Wait()
}
