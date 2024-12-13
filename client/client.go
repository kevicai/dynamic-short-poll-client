package client

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
	ApiUrl string
}

func NewJobClient() *JobClient {
	LoadStats()

	return &JobClient{}
}

var (
	infrequentPollInterval time.Duration = 1 * time.Second // poll interval outside of the 95% data interval
)

func (c *JobClient) apiUrl() string {
	if c.ApiUrl == "" {
		return DefaultBaseURL
	}
	return c.ApiUrl
}

func (c *JobClient) CreateJob(onComplete func(job *Job)) (jobID string, err error) {
	url := fmt.Sprintf("%s/job", c.apiUrl())

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
	url := fmt.Sprintf("%s/status?id=%s", c.apiUrl(), jobID)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("server response status is not OK")
	}

	var status JobStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return "", err
	}
	return status.Result, nil
}

// Poll for status automatically
func (c *JobClient) shortPollStatus(job Job, onComplete func(job *Job)) {
	startTime := time.Now()

	frequentPollInterval := adjustPollInterval()

	pollCount := 0
	for {
		timeElapsed := time.Since(startTime)

		// poll more frequently when time is within 95% data range
		if timeElapsed >= stats.AvgTime-2*stats.StdDeviation && timeElapsed <= stats.AvgTime+2*stats.StdDeviation {
			time.Sleep(frequentPollInterval)
		} else {
			time.Sleep(infrequentPollInterval) // poll infrequently
		}
		pollCount++

		status, err := c.queryJobStatus(job.ID)
		if err != nil {
			continue
		}

		// job completed
		if status != "pending" {
			endTime := time.Now()
			timeDiff := endTime.Sub(startTime)
			// fmt.Printf("Times Polled: %v. Time taken: %v\n", pollCount, timeDiff)

			// if error status are rare and has different job durations, can also add a check to only update stats for completed jobs
			go updateStats(timeDiff)

			job.Status = status
			go onComplete(&job)
			return
		}
	}
}

func adjustPollInterval() time.Duration {
	if stats.NumJobs == 0 {
		return infrequentPollInterval // default
	}

	// dynamically adjust the infrequent poll interval if its too large
	if infrequentPollInterval > stats.AvgTime-2*stats.StdDeviation {
		infrequentPollInterval /= 2
	}
	// dynamically adjust poll interval during infrequent times
	infrequentPollInterval = (stats.AvgTime - 2*stats.StdDeviation) / 4

	// poll 10 times during the frequent distribution of data
	frequentPollInterval := 4 * stats.StdDeviation / 10
	frequentPollInterval = max(frequentPollInterval, time.Second/2) // cap poll interval at 0.5 seconds

	// fmt.Printf("Infrequent interval: %v\n", infrequentPollInterval)
	// fmt.Printf("Frequent interval: %v\n", frequentPollInterval)
	return frequentPollInterval
}
