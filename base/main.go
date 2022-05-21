package main

import (
	"encoding/json"
	"math/rand"
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
			register: RegisterService{
				service: Service{
					// randome id
					Id: rand.Intn(99999999999),
				},
				backChan: make(chan RegisterServiceResult),
			},
			getServices: GetServices{
				backChan: make(chan []Service),
			},
		}

		jobs <- job
		<-job.register.backChan
		result := <-job.getServices.backChan

		// retunr the result2
		json.NewEncoder(w).Encode(result)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		register(w, r, jobs)
	})

	log.Print("Application is ready to serve requests.")
	http.ListenAndServe(":80", nil)
}
