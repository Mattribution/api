package functions

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/mattribution/api/internal/app"
	"github.com/mattribution/api/internal/pkg/postgres"
)

var (
	mux    *http.ServeMux
	tracks *postgres.Tracks
	gif    = []byte{
		71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 0, 0, 0,
		255, 255, 255, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0,
		1, 0, 1, 0, 0, 2, 1, 68, 0, 59,
	}
)

func init() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	CSQLInstanceConnectionName := os.Getenv("CSQL_INSTANCE_CONNECTION_NAME")

	// Setup router
	mux = http.NewServeMux()
	mux.HandleFunc("/tracks/new", newTrack)

	// Setup db connection
	db, err := postgres.NewCloudSQLClient(dbUser, dbPass, dbName, CSQLInstanceConnectionName)
	if err != nil {
		panic(err)
	}
	// Only allow 1 connection to the database to avoid overloading
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	// Setup services
	tracks = &postgres.Tracks{
		DB: db,
	}
}

// FunctionsEntrypoint represents cloud function entry point
func FunctionsEntrypoint(w http.ResponseWriter, r *http.Request) {
	mux.ServeHTTP(w, r)
}

func newTrack(w http.ResponseWriter, r *http.Request) {
	// Get pixel data from client
	v := r.URL.Query()
	rawEvent := v.Get("data")
	data, err := base64.StdEncoding.DecodeString(rawEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Unmarshal
	track := app.Track{}
	if err := json.Unmarshal(data, &track); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Set received at
	now := time.Now()
	track.ReceivedAt = &now

	// Grab IP
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	track.IP = &ip

	// Store raw track
	// TODO: Make this all a transaction to fail if one can or can't
	newTrackID, err := tracks.Store(track)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	track.ID = newTrackID

	// Write gif back to client
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Content-Type", "image/gif")
	w.Write(gif)
}
