package main

import (
	"encoding/json"
	"net/http"
	"os"

	"base/Registry"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var jobs = make(chan registry.PoolJob, 10)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Print("Starting the application...")

	go registry.PoolWorker(jobs)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		job := registry.PoolJob{
			GetServices: registry.GetServices{
				BackChan: make(chan []registry.Service),
			},
		}
		jobs <- job
		result := <-job.GetServices.BackChan

		json.NewEncoder(w).Encode(result)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		register(w, r, jobs)
	})

	log.Print("Application is ready to serve requests.")
	http.ListenAndServe(":80", nil)
}
