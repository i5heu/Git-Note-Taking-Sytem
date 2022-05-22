package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func register(w http.ResponseWriter, r *http.Request, jobs chan PoolJob) {
	// time start
	start := time.Now()

	var service Service
	err := json.NewDecoder(r.Body).Decode(&service)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	job := PoolJob{
		register: RegisterService{
			service: Service{
				Id:   service.Id,
				Name: service.Name,
			},
			backChan: make(chan RegisterServiceResult),
		},
	}

	jobs <- job
	result := <-job.register.backChan

	if result.results == "OK" {
		success(w)
	} else {
		http.Error(w, result.results, http.StatusBadRequest)
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
