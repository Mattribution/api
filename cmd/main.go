package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/mattribution/api/pkg/http"
	"github.com/mattribution/api/pkg/postgres"
)

const (
	host     = "localhost"
	port     = 5432
	username = "postgres"
	password = "password"
	dbName   = "mattribution"
)

func main() {

	// Make connection to DB
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbName)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	// Setup
	trackService := postgres.TrackService{db}
	kpiService := postgres.KPIService{db}
	conversionService := postgres.ConversionService{db}
	billingEventService := postgres.BillingEventService{db}
	campaignService := postgres.CampaignService{db}

	httpHandler := http.NewHandler(trackService, kpiService, conversionService, billingEventService, campaignService)
	log.Fatal(httpHandler.Serve(":3001"))

}
