package main

import (
	"encoding/json"
	"net/http"
	"os"

	channels "base/Channels"
	gitManager "base/GitManager"
	queue "base/Queue"
	registry "base/Registry"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Print("Starting the application...")

	go queue.QueueWorker(channels.QueueChan)
	go registry.PoolWorker(channels.RegistryChan)
	go gitManager.GitManager()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		job := registry.PoolJob{
			GetServices: registry.GetServices{
				BackChan: make(chan []registry.Service),
			},
		}
		channels.RegistryChan <- job
		result := <-job.GetServices.BackChan

		json.NewEncoder(w).Encode(result)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		register(w, r, channels.RegistryChan)
	})

	log.Print("Application is ready to serve requests.")
	http.ListenAndServe(":80", nil)
}
