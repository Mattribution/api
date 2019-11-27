package main

import (
	"log"

	"github.com/mattribution/api/pkg/http"
	"github.com/mattribution/api/pkg/postgres"
)

func main() {

	// Setup
	trackService, err := postgres.NewTrackService("localhost", "postgres", "password", "mattribution", 5432)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
	}
	kpiService, err := postgres.NewKPIService("localhost", "postgres", "password", "mattribution", 5432)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
	}

	httpHandler := http.NewHandler(trackService, kpiService)
	log.Fatal(httpHandler.Serve(":3001"))

}
