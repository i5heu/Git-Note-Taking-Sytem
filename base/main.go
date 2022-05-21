package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var jobs = make(chan RegistryJob, 10)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Print("Starting the application...")

	go PoolWorker(jobs)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		job := RegistryJob{
			id:       1,
			randomno: 12,
			backChan: make(chan RegistryResult),
		}

		jobs <- job
		result := <-job.backChan

		w.Write([]byte(fmt.Sprintf("%d\n", result.results)))
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		register(w, r, jobs)
	})

	log.Print("Application is ready to serve requests.")
	http.ListenAndServe(":80", nil)
}
