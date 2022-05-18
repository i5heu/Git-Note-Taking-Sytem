package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var pool = sync.Pool{
	New: func() interface{} {
		return &http.Request{}
	},
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Print("Starting the application...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello, world! This is Tyche")
		w.Write([]byte("Hello, world!"))
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		register(w, r)
	})

	log.Print("Application is ready to serve requests.")
	http.ListenAndServe(":80", nil)
}
