package queue

import (
	registry "base/Registry"
	"time"

	"github.com/rs/zerolog/log"
)

type QueueJob struct {
	QueueItem     QueueItem
	Register      QueueItem
	GetStatus     GetStatus
	GetAllStatus  GetAllStatus
	FileAuthority FileAuthority

	queueWorkerTick queueWorkerTick
}

//enum status
const (
	WAITING int = iota
	RUNNING
	FINISHED
)

type QueueItem struct {
	Service       registry.Service
	Priority      int // if -1 the job will be dropped
	RunAfterwards []registry.Service
	SlotOpen      chan bool

	status                int
	initFromRunAfterwards bool
}

type FileAuthority struct {
	Service                       registry.Service
	FileAuthorityGrant            chan bool // if chanel is given a File Authority will be granted
	FileAuthorityGrantTermination chan bool
	FileAuthorityTimeOutDuration  time.Duration
}

type GetStatus struct {
	QueueItem QueueItem
	BackChan  chan int
}

type GetAllStatus struct {
	Service  registry.Service
	BackChan chan []GetStatusItems
}

type GetStatusItems struct {
	Service registry.Service
	Status  int
}

type setStatus struct {
	QueueItem QueueItem
	Status    int
}

type queueWorkerTick struct {
	marker int
}

func QueueWorker(jobs chan QueueJob) {
	var queues []QueueItem

	go heartBeat(jobs)

	for job := range jobs {

		if job.Register.SlotOpen != nil {
			queues = append(queues, job.Register)
		}

		if job.GetStatus.BackChan != nil {
			job.GetStatus.BackChan <- job.GetStatus.QueueItem.status
		}

		if job.GetAllStatus.BackChan != nil {
			HandleGetStatus(queues, job)
		}

		if job.queueWorkerTick.marker > -1 {
			log.Debug().Int("marker", job.queueWorkerTick.marker).Msg("QueueWorker: tick")
		}

		if job.FileAuthority.FileAuthorityGrant != nil {
			job.FileAuthority.Service.HasFileAuthority = true
			job.FileAuthority.FileAuthorityGrant <- true
			go func() {
				time.Sleep(job.FileAuthority.FileAuthorityTimeOutDuration)
				job.FileAuthority.Service.HasFileAuthority = false

				//TODO wait for all files to be closed
				job.FileAuthority.FileAuthorityGrant <- false
			}()
			//sketchy? GC will remove?
			go func() {
				<-job.FileAuthority.FileAuthorityGrantTermination
				job.FileAuthority.Service.HasFileAuthority = false
				
				//TODO wait for all files to be closed
				job.FileAuthority.FileAuthorityGrant <- false
			}()
		}

	}
}

func heartBeat(jobs chan QueueJob) {
	tick := 0

	for range time.Tick(time.Second * 1) {
		jobs <- QueueJob{queueWorkerTick: queueWorkerTick{tick}}
		tick++
	}
}

func HandleGetStatus(queues []QueueItem, job QueueJob) {
	var statusItems []GetStatusItems

	for _, item := range queues {
		if item.Service.Id == job.GetAllStatus.Service.Id {
			statusItems = append(statusItems, GetStatusItems{Service: item.Service, Status: item.status})
		}
	}

	job.GetAllStatus.BackChan <- statusItems
}
