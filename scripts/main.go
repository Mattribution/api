package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	wr "github.com/mroth/weightedrand"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // always seed random!

	var host, port, username, password, dbName = "localhost", "postgres", "password", "mattribution", 5432
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbName)
	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	randMedium := wr.NewChooser(
		wr.Choice{Item: "Google AdWords", Weight: 8},
		wr.Choice{Item: "Search", Weight: 4},
		wr.Choice{Item: "Blog", Weight: 2},
		wr.Choice{Item: "Twitter", Weight: 1},
	)

	for i := 0; i < 1000; i++ {
		medium := randMedium.Pick().(string)
	}
}
