# Dynamic Short Poll Client

A client library in Golang that checks the status of simulated video translation jobs in a simple server.

The simple server in `/server` has two endpoints:
- 'POST' /job: creates a job and returns a job object. The job completes after a random configurable delay. Currently set at 10-20s
- 'GET' /status: Returns the status of the video translation jobs as a JSON:
 
    {“result”: “pending” or “error” or “completed}

## Details on the dynamic short polling implementation

For the client to know the status of a job, ideally, a web socket should be used to signal the client whenever a job completes or fails. But in this case, the server only accepts a GET request to check the status of the job, short polling is the only viable option to check the status regularly. 

We could let users check the status of the job by repeatedly calling a check status method by themselves. However, this runs into the risk of calling the server too frequently, which uses up lots of resources, or infrequently, which causes delayed results. To address this, the client automatically checks the status of the job and notifies the user with a callback function after the client has confirmed that the job is finished. 

To check the job at regular intervals, the client currently assumes that the job takes at least a few seconds to complete and that the job completion time approximately follows a normal distribution. The client first sets in periodic checking intervals to 1 second, then after each completed job, the client updates its internal stats to keep track of the average and standard deviation of the job durations. This makes it possible for the client to short poll using dynamic intervals, which polls the status infrequently when the job execution time is far from the normal, and more frequently when the job execution time falls within 2 standard deviations from the normal. This ensures that fewer polls are performed overall, but still maintains a low delay as to calling it on a fixed interval.

## Experiment

The following experiment shows the performance before and after training on 10 jobs with duration between 10-20 seconds.

### Before training:

![alt text](/imgs/pre-train.png)

### After training:

![alt text](/imgs/post-train.png)

There is a noticable improvement of poll times from the dynamic poll interval implementation after training. Before, a job with around 11 seconds took 11 polls, where as after training it only took 6 polls check the job status. 

## Installation 

The client library can be installed with go 1.23 or higher:

```
go get github.com/kevicai/job-status-api/client@latest
```

## Usage 

```go
import "github.com/kevicai/job-status-api/client"

clnt := client.NewJobClient()

// A callback function for how to handle a job after its completion
handleComplete := func(job *client.Job) {
    fmt.Println("Job completed:", job.ID, "| Status:", job.Status)
}

// To create a job with the callback function
jobID, err := clnt.CreateJob(handleComplete)
```

## Integration test

To spin up the server and test the client:

```bash
make test
```

Server logs are output to `server.log` file. 

To remove the log file and stop the server, run:

```
make clean
```

## Next steps:

- The client should save its stats periodically, eg. after 10 jobs to a binary file, and load the binary file each time it's initialized. This way the client does not need to be retrained, and the stats can be more accurate. 
- Translation jobs should be categorized based on some observable pattern to better predict the duration. For example, based on the length of the text, or the server capacity. This would require a more complex model.
