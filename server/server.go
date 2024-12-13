package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kevicai/job-status-api/server/job"
)

var (
	jobManager = job.NewJobManager()
)

func jobStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jobID := r.URL.Query().Get("id")
	if jobID == "" {
		http.Error(w, `{"error": "id is required"}`, http.StatusBadRequest)
		return
	}

	status, ok := jobManager.GetJobStatus(jobID)
	if !ok {
		http.Error(w, `{"error": "job not found"}`, http.StatusNotFound)
		return
	}

	// send a JSON with only the status
	response := map[string]string{"result": status.String()}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "failed to encode status JSON"}`, http.StatusInternalServerError)
		return
	}
}

func createJobHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	job := jobManager.CreateJob()
	w.WriteHeader(http.StatusCreated)

	// send the created job as a JSON
	json.NewEncoder(w).Encode(job)
}

func main() {
	http.HandleFunc("GET /status", jobStatusHandler)
	http.HandleFunc("POST /job", createJobHandler)

	log.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
