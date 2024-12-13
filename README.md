# Dynamic Short Poll Client

A client library in Golang to simulate the video translation jobs in the server.

This client calls a simple server in `/server` which has two end points:
- 'POST' /job: creates a job and returns a job object. The job completes after a random configurable delay. Currently set at 10-20s
- 'GET' /status: Returns the status of the video translation jobs as a JSON:
 
    {“result”: “pending” or “error” or “completed}

## Installation 

The client library can be installed with go 1.23 or higher:

```
go install github.com/kevicai/job-status-api/client@latest
```

## Usage 

```go
import "github.com/kevicai/job-status-api/client"

client := client.NewJobClient()
```

## Deliverable

A public git repository with your code 
Write a small integration test that spins up your server and uses your client library to demonstrate the usage and print logs.
Write a small doc on how to use your client library.

![alt text](/imgs/pre-train.png)

![alt text](/imgs/post-train.png)