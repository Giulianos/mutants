package main

import (
	"github.com/Giulianos/mutants/internal/analyzer"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

const (
	NSQD_HOST="nsqd:4151"
)

func main() {
	// read env configuration
	nsqdHost, defined := os.LookupEnv("NSQD_HOST")
	if !defined {
		nsqdHost = NSQD_HOST
	}

	// create deps
	eventPublisher := analyzer.NewEventPublisher(nsqdHost)

	r := mux.NewRouter()
	r.Handle("/mutant", analyzer.NewController(eventPublisher)).Methods("POST")

	s := http.Server{
		Addr:         ":80",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}
