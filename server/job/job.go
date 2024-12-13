package job

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// Use status as an enum for smaller size
type JobStatus int

const (
	Pending JobStatus = iota // 0
	Completed
	Error
)

// Returns the string representation of the JobStatus enum
func (js JobStatus) String() string {
	return [...]string{"pending", "completed", "error"}[js]
}

type Job struct {
	ID     string    `json:"id"`
	Status JobStatus `json:"status"`
}

// Customizes Job's JSON representation
func (job Job) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":     job.ID,
		"status": job.Status.String(),
	})
}

func (job *Job) Start(onComplete func(job *Job)) {
	executionTime := time.Duration(minJobTime + rand.Float64()*(maxJobTime-minJobTime))

	// simulate processing delay
	time.Sleep(executionTime * time.Second)

	// flip a coin to decide whether the job fails or completes
	switch rand.Intn(2) {
	case 0:
		job.Fail()
	default:
		job.Complete()
	}

	fmt.Printf("Time taken: %v\n", executionTime*time.Second)
	onComplete(job)
}

func (job *Job) Fail() {
	fmt.Printf("Job: %s [Failed]\n", job.ID)
	job.Status = Error
}

func (job *Job) Complete() {
	fmt.Printf("Job: %s [Completed]\n", job.ID)
	job.Status = Completed
}
