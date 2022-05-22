package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var jobs = make(chan PoolJob, 10)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Print("Starting the application...")

	go PoolWorker(jobs)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		job := PoolJob{
			getServices: GetServices{
				backChan: make(chan []Service),
			},
		}
		jobs <- job
		result := <-job.getServices.backChan

		json.NewEncoder(w).Encode(result)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		register(w, r, jobs)
	})

	log.Print("Application is ready to serve requests.")
	http.ListenAndServe(":80", nil)
}
