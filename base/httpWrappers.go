package main

import (
	"encoding/json"
	"net/http"
	"time"

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

func success(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
