package client

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"sync"
	"time"
)

// Stats struct to hold all statistical data
type PollStats struct {
	NumJobs      int           // Number of jobs completed
	AvgTime      time.Duration // Average time of jobs
	StdDeviation time.Duration // Standard deviation of duration
	VarianceSum  float64       // Variance sum of duration
}

type PollStatsBinary struct {
	NumJobs      int32
	AvgTime      int64
	StdDeviation int64
	VarianceSum  int64
}

const statsFile = ".poll-stats"

var (
	stats      PollStats
	statsMutex sync.Mutex
)

func updateStats(timeTaken time.Duration) {
	statsMutex.Lock()
	defer statsMutex.Unlock()

	stats.NumJobs++
	totalTime := stats.AvgTime*time.Duration(stats.NumJobs-1) + timeTaken
	stats.AvgTime = totalTime / time.Duration(stats.NumJobs)

	// update standard deviation
	if stats.NumJobs > 1 {
		diff := float64(timeTaken-stats.AvgTime) / float64(time.Second) // divide by seconds for a smaller value
		stats.VarianceSum += diff * diff
		stats.StdDeviation = time.Duration(math.Sqrt(stats.VarianceSum/float64(stats.NumJobs)) * float64(time.Second))
		// fmt.Printf("New stats: average time=%v, std deviation=%v\n", stats.AvgTime, stats.StdDeviation)
	}

	// buffer write for every 10 new stats
	if stats.NumJobs%10 == 0 {
		SaveStats()
	}
}

// LoadStats reads stats from the binary file
func LoadStats() error {
	statsMutex.Lock()
	defer statsMutex.Unlock()

	file, err := os.Open(statsFile)
	if err != nil {
		if os.IsNotExist(err) {
			// file doesn't exist, start fresh
			stats = PollStats{
				NumJobs:      0,
				AvgTime:      time.Duration(0),
				StdDeviation: time.Duration(0),
				VarianceSum:  0,
			}
			return nil
		}
		return err
	}
	defer file.Close()

	// read the binary metadata file
	binaryStats := PollStatsBinary{}
	if err := binary.Read(file, binary.LittleEndian, &binaryStats); err != nil {
		return err
	}

	stats = PollStats{
		NumJobs:      int(binaryStats.NumJobs),
		AvgTime:      time.Duration(binaryStats.AvgTime),
		StdDeviation: time.Duration(binaryStats.StdDeviation) * time.Millisecond,
		VarianceSum:  math.Float64frombits(uint64(binaryStats.VarianceSum)),
	}

	return nil
}

// Saves stats to the binary file
func SaveStats() error {
	fmt.Println("Saving stats")
	statsMutex.Lock()
	defer statsMutex.Unlock()

	file, err := os.Create(statsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	binaryStats := PollStatsBinary{
		NumJobs:      int32(stats.NumJobs),
		AvgTime:      int64(stats.AvgTime),
		StdDeviation: int64(stats.StdDeviation.Milliseconds()),
		VarianceSum:  int64(math.Float64bits(stats.VarianceSum)),
	}

	if err := binary.Write(file, binary.LittleEndian, &binaryStats); err != nil {
		return err
	}

	return nil
}
