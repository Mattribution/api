package functions

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/mattribution/api/internal/app"
	"github.com/mattribution/api/internal/pkg/postgres"
)

type Handler struct {
	Tracks app.Tracks
}

const (
	newTrackPath               = "/tracks/new"
	invalidRequestError        = "The request you sent is invalid. Please reformat the request and try again."
	invalidBase64EncodingError = "The data sent was not Base64 encoded. Please encode the data and try again."
)

var (
	gif = []byte{
		71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 0, 0, 0,
		255, 255, 255, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0,
		1, 0, 1, 0, 0, 2, 1, 68, 0, 59,
	}
	mux     *http.ServeMux
	handler *Handler
)

func init() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	CSQLInstanceConnectionName := os.Getenv("CSQL_INSTANCE_CONNECTION_NAME")

	// Setup db connection
	db, err := postgres.NewCloudSQLClient(dbUser, dbPass, dbName, CSQLInstanceConnectionName)
	if err != nil {
		panic(err)
	}
	// Only allow 1 connection to the database to avoid overloading
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	// Setup services
	handler = &Handler{
		Tracks: &postgres.Tracks{
			DB: db,
		},
	}

	// Setup router
	mux = http.NewServeMux()
	mux.HandleFunc(newTrackPath, handler.newTrack)
}

// FunctionsEntrypoint represents cloud function entry point
func FunctionsEntrypoint(w http.ResponseWriter, r *http.Request) {
	mux.ServeHTTP(w, r)
}

// ~=~=~=~=~=~=~=~=
// Tracks
// ~=~=~=~=~=~=~=~=

func (h *Handler) newTrack(w http.ResponseWriter, r *http.Request) {
	// Get pixel data from client
	v := r.URL.Query()
	rawEvent := v.Get("data")
	data, err := base64.StdEncoding.DecodeString(rawEvent)
	if err != nil {
		http.Error(w, invalidBase64EncodingError, http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Unmarshal
	track := app.Track{}
	if err := json.Unmarshal(data, &track); err != nil {
		http.Error(w, invalidRequestError, http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Grab IP
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	track.IP = ip

	// Store raw track
	newTrackID, err := h.Tracks.Store(track)
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

// ~=~=~=~=~=~=~=~=
// Kpis
// ~=~=~=~=~=~=~=~=
