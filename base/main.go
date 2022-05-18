package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Print("Starting the application...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello, world! This is Tyche")
		w.Write([]byte("Hello, world!"))
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		
		fmt.Println("register:")
		w.Write([]byte("Hello, world!"))
	})

	log.Print("Application is ready to serve requests.")
	http.ListenAndServe(":80", nil)
}
