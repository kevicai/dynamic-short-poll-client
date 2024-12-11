package job

import (
	"github.com/google/uuid"
)

const (
	minJobTime float64 = 10 // in seconds
	maxJobTime float64 = 30
)

type JobManager struct {
	jobStatuses TypedSyncMap // stores job statuses in memory
}

// Creates a new JobManager
func NewJobManager() *JobManager {
	return &JobManager{}
}

func (jm *JobManager) CreateJob() Job {
	jobID := uuid.New().String()
	job := Job{
		ID:     jobID,
		Status: Pending,
	}
	jm.AddJobStatus(job.ID, job.Status)

	go job.Start(jm.HandleJobComplete)

	return job
}

func (jm *JobManager) GetJobStatus(jobID string) (status JobStatus, ok bool) {
	status, ok = jm.jobStatuses.Load(jobID)
	return
}

func (jm *JobManager) AddJobStatus(jobID string, status JobStatus) {
	jm.jobStatuses.Store(jobID, status)
}

func (jm *JobManager) HandleJobComplete(job *Job) {
	// should remove and put into storage in real world for completed jobs
	jm.UpdateJobStatus(job.ID, job.Status)
	// jm.RemoveJob(jobID)
}

func (jm *JobManager) UpdateJobStatus(jobID string, newStatus JobStatus) {
	jm.jobStatuses.Store(jobID, newStatus)
}

func (jm *JobManager) RemoveJob(jobID string) {
	jm.jobStatuses.Delete(jobID)
}
