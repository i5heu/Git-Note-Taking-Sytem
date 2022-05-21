package main

type PoolJob struct {
	register    RegisterService
	getServices GetServices
}
type Service struct {
	Id int `json:"id"`
}
type RegisterService struct {
	service  Service
	backChan chan RegisterServiceResult
}
type RegisterServiceResult struct {
	results string
}
type GetServices struct {
	backChan chan []Service
}

func PoolWorker(jobs chan PoolJob) {
	var servicePoll []Service

	for job := range jobs {

		if job.register.backChan != nil {
			servicePoll = append(servicePoll, job.register.service)
			job.register.backChan <- RegisterServiceResult{results: "OK"}
		}

		if job.getServices.backChan != nil {
			job.getServices.backChan <- servicePoll
		}
	}
}
