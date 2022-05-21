package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func register(w http.ResponseWriter, r *http.Request, jobs chan RegistryJob) {
	// time start
	start := time.Now()

	var pr PluginRegister
	err := json.NewDecoder(r.Body).Decode(&pr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	job := RegistryJob{
		id:       1,
		randomno: 12,
		backChan: make(chan RegistryResult),
	}

	jobs <- job
	result := <-job.backChan

	w.Write([]byte(fmt.Sprintf("%d\n", result.results)))

	// success(w)

	// time end
	end := time.Now()
	elapsed := end.Sub(start)
	log.Debug().Timestamp().Int("elapsed-micro", int(elapsed.Microseconds())).Str("API-Method", "register").Msg("Time elapsed")
}

func success(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
