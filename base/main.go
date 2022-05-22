package main

import (
	"encoding/json"
	"net/http"
	"os"

	queue "base/Queue"
	registry "base/Registry"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var registryChan = make(chan registry.PoolJob, 50)
var queueChan = make(chan queue.QueueJob, 50)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Print("Starting the application...")

	go queue.QueueWorker(queueChan)
	go registry.PoolWorker(registryChan)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		job := registry.PoolJob{
			GetServices: registry.GetServices{
				BackChan: make(chan []registry.Service),
			},
		}
		registryChan <- job
		result := <-job.GetServices.BackChan

		json.NewEncoder(w).Encode(result)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		register(w, r, registryChan)
	})

	log.Print("Application is ready to serve requests.")
	http.ListenAndServe(":80", nil)
}
