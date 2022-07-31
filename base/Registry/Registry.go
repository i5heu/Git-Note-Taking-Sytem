package registry

import (
	"github.com/google/uuid"
)

type PoolJob struct {
	Register    RegisterService
	GetServices GetServices
}
type Service struct {
	Id               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	HasFileAuthority bool      `json:"hasFileAuthority"`
}
type RegisterService struct {
	Service  Service
	BackChan chan RegisterServiceResult
}
type RegisterServiceResult struct {
	Results string
}
type GetServices struct {
	BackChan chan []Service
}

func PoolWorker(jobs chan PoolJob) {
	var servicePoll []Service

	for job := range jobs {

		if job.Register.BackChan != nil {
			servicePoll = append(servicePoll, job.Register.Service)
			job.Register.BackChan <- RegisterServiceResult{Results: "OK"}
		}

		if job.GetServices.BackChan != nil {
			job.GetServices.BackChan <- servicePoll
		}
	}
}
