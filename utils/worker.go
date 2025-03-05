package utils

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID   int
	Size int
	Task func()
}

func NewWorker(id int, jobs <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		fmt.Printf("Worker %d processing job %d (Size: %dMB)\n", id, job.ID, job.Size)
		job.Task()
	}

}

func DispatchJobs(jobs []Job, maxWorkers, memoryLimit int) {
	jobQueue := make(chan Job, len(jobs))
	var wg sync.WaitGroup

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go NewWorker(i+1, jobQueue, &wg)
	}
	currentMemory := 0

	for _, job := range jobs {
		for currentMemory+job.Size > memoryLimit {
			fmt.Println("Memory limit reached, waiting...")
			time.Sleep(500 * time.Millisecond) // Simulate waiting
		}
		currentMemory += job.Size
		go func(j Job) {
			jobQueue <- j
			time.AfterFunc(1*time.Second, func() {
				currentMemory -= j.Size
			})
		}(job)
	}

	close(jobQueue)
	wg.Wait()
}
