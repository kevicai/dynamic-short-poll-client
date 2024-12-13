package translation

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const DefaultBaseURL = "http://localhost:8080"

// The translation would also be stored in here
type Job struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type JobStatus struct {
	Result string `json:"result"`
}

// JobClient struct can be constructed explicitly
type JobClient struct {
	BaseURL      string
	PollInterval time.Duration
}

func (c *JobClient) baseUrl() string {
	if c.BaseURL == "" {
		return DefaultBaseURL
	}
	return c.BaseURL
}

func (c *JobClient) pollInterval() time.Duration {
	return 1 * time.Second
}

func (c *JobClient) CreateJob(onComplete func(job *Job)) (jobID string, err error) {
	url := fmt.Sprintf("%s/job", c.baseUrl())

	res, err := http.Post(url, "application/json", nil)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return "", errors.New("job creation failed")
	}

	var createdJob Job
	if err := json.NewDecoder(res.Body).Decode(&createdJob); err != nil {
		return "", err
	}

	go c.shortPollStatus(createdJob, onComplete)

	return createdJob.ID, nil
}

// Calls the server to get the latest job status
func (c *JobClient) queryJobStatus(jobID string) (string, error) {
	url := fmt.Sprintf("%s/status?id=%s", c.baseUrl(), jobID)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("response status not OK")
	}

	var status JobStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return "", err
	}
	return status.Result, nil
}

// Poll for status automatically
func (c *JobClient) shortPollStatus(job Job, onComplete func(job *Job)) {
	for {
		time.Sleep(c.pollInterval())

		status, err := c.queryJobStatus(job.ID)
		if err != nil {
			continue
		}

		// job completed
		if status != "pending" {
			job.Status = status
			go onComplete(&job)
			return
		}
	}
}
