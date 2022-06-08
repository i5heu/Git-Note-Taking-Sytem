package main

import (
	"encoding/json"
	"net/http"
	"time"

	channels "base/Channels"
	files "base/Files"
	"base/Registry"

	"github.com/rs/zerolog/log"
)

func register(w http.ResponseWriter, r *http.Request, jobs chan registry.PoolJob) {
	// time start
	start := time.Now()

	var service registry.Service
	err := json.NewDecoder(r.Body).Decode(&service)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	job := registry.PoolJob{
		Register: registry.RegisterService{
			Service: registry.Service{
				Id:   service.Id,
				Name: service.Name,
			},
			BackChan: make(chan registry.RegisterServiceResult),
		},
	}

	jobs <- job
	result := <-job.Register.BackChan

	if result.Results == "OK" {
		success(w)
	} else {
		http.Error(w, result.Results, http.StatusBadRequest)
	}

	// time end
	end := time.Now()
	elapsed := end.Sub(start)
	log.Debug().Timestamp().Int("elapsed-micro", int(elapsed.Microseconds())).Str("API-Method", "register").Msg("Time elapsed")
}

func requestsTaskSlot(w http.ResponseWriter, r *http.Request) {
	log.Debug().Str("API-Method", "requestsTaskSlot").Msg("TaskSlot")
}

type getFileRequest struct {
	FilePath string
	Lock     bool
	BackURL  string
}

func getFile(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// var requestData getFileRequest
	// err := json.NewDecoder(r.Body).Decode(&requestData)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	job := files.FileJob{
		GetAndLockFile: files.GetFile{
			FilePath: "/today.txt",
			Lock:     true,
			BackChan: make(chan files.File),
		},
		TerminationChan: make(chan bool),
	}

	channels.FileChan <- job
	if <-job.TerminationChan {
		http.Error(w, "File is locked", http.StatusLocked)
		return
	}

	result := <-job.GetAndLockFile.BackChan

	json.NewEncoder(w).Encode(
		files.File{
			FilePath:   result.FilePath,
			LockStatus: result.LockStatus,
			LockExpire: result.LockExpire,
			Content:    result.GetContent(),
		})

	end := time.Now()
	elapsed := end.Sub(start)
	log.Debug().Timestamp().Int("elapsed-micro", int(elapsed.Microseconds())).Str("API-Method", "getFile").Msg("Time elapsed")
}

func success(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
