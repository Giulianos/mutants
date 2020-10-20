package main

import (
	"github.com/Giulianos/mutants/internal/stats"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	MONGO_HOST = "mongodb://localhost:27017"
	MONGO_DB   = "stats"
	NSQ_LOOKUP = "nsqlookupd:4161"
)

func main() {
	// read env configuration
	mongoHost, defined := os.LookupEnv("MONGO_HOST")
	if !defined {
		mongoHost = MONGO_HOST
	}
	mongoDB, defined := os.LookupEnv("MONGO_DB")
	if !defined {
		mongoDB = MONGO_DB
	}
	nsqLookup, defined := os.LookupEnv("NSQ_LOOKUP")
	if !defined {
		nsqLookup = NSQ_LOOKUP
	}

	// create deps (consider using a DI framework if adding more deps)
	repo, err := stats.NewMongoRepository(mongoHost, mongoDB)
	if err != nil {
		log.Fatal("unable to create repository:", err)
	}
	service := stats.NewService(repo)
	eventHandler := stats.NewEventHandler(service)

	// define routes
	r := mux.NewRouter()
	r.Handle("/stats", stats.NewController(service)).Methods("GET")
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	// create http server
	s := http.Server{
		Addr:         ":80",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// start listening events
	done := make(chan struct{})
	go stats.StartListener(nsqLookup, &eventHandler, done)

	// start listening requests
	log.Fatal(s.ListenAndServe())
}
