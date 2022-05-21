package main

type RegistryJob struct {
	id       int
	randomno int
	backChan chan RegistryResult
}
type RegistryResult struct {
	results int
}

func PoolWorker(jobs chan RegistryJob) {
	counter := 0
	for job := range jobs {
		output := RegistryResult{counter}
		job.backChan <- output
		counter++
	}
}
