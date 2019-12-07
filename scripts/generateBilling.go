package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mattribution/api/pkg/api"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // always seed random!

	var host, username, password, dbName = "localhost", "postgres", "password", "mattribution"
	port := 5432
	// Amount of users to mock
	userCount := 100
	// Mock owner id
	var ownerID = 1
	bar := pb.StartNew(userCount)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbName)
	log.Println(psqlInfo)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	for userID := 0; userID < userCount; userID++ {
		bar.Increment()

		shouldCharge := rand.Intn(2)
		if shouldCharge == 0 {
			randAmount := rand.Float32() * 3000
			storeBillingEvent(db, api.BillingEvent{
				OwnerID: ownerID,
				UserID:  userID,
				Amount:  randAmount,
			})
		}
	}

	bar.Increment()
}

func storeBillingEvent(db *sqlx.DB, be api.BillingEvent) {
	sqlStatement :=
		`INSERT INTO public.billing_events (owner_id, user_id, amount, created_at)
	VALUES($1, $2, $3, $4)`

	_, err := db.Exec(sqlStatement, be.OwnerID, be.UserID, be.Amount, time.Now().Format(time.RFC3339))
	if err != nil {
		log.Printf("ERROR: %s", err)
	}
}
