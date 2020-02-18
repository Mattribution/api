package functions

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/mattribution/api/internal/app"
	"github.com/mattribution/api/internal/pkg/postgres"
)

type Handler struct {
	Tracks app.Tracks
	Kpis   app.Kpis
}

type ContextKey string

const (
	invalidRequestError                   = "The request you sent is invalid. Please reformat the request and try again."
	invalidBase64EncodingError            = "The data sent was not Base64 encoded. Please encode the data and try again."
	internalError                         = "We experienced an internal error. Please try again later."
	mockOwnerID                int64      = 0
	ContextKeyOwnerID          ContextKey = "ownerID"
)

var (
	gif = []byte{
		71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 0, 0, 0,
		255, 255, 255, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0,
		1, 0, 1, 0, 0, 2, 1, 68, 0, 59,
	}
	router  *mux.Router
	handler *Handler
)

func init() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	// Setup db connection
	db, err := postgres.NewCloudSQLClient(dbUser, dbPass, dbName, dbHost)
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
		Kpis: &postgres.Kpis{
			DB: db,
		},
	}

	// Setup router
	router = handler.Router()
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextKeyOwnerID, mockOwnerID)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Router generates the routes
func (h *Handler) Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/tracks/new", h.newTrack).Methods("GET")
	router.HandleFunc("/kpis", h.newKpi).Methods("POST")
	router.HandleFunc("/kpis/{id:[0-9]+}", h.deleteKpi).Methods("DELETE")
	router.HandleFunc("/kpis", h.listKpis).Methods("GET")
	router.Use(h.AuthMiddleware)
	return router
}

// FunctionsEntrypoint represents cloud function entry point
func FunctionsEntrypoint(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
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
		log.Println("Error storing track: ", err)
		return
	}
	track.ID = newTrackID

	// Write gif back to client
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Content-Type", "image/gif")
	w.Write(gif)
}

// func (h *Handler) getTrackJourneyAggregate(w http.ResponseWriter, r *http.Request) {
// 	q := r.URL.Query()
// 	columnName := q.Get("column_name")

// 	ownerIDInterface := r.Context().Value(ContextKeyOwnerID)
// 	ownerID, ok := ownerIDInterface.(int64)
// 	if !ok {
// 		http.Error(w, "owner id error", http.StatusBadRequest)
// 		return
// 	}
// 	// Get aggregate data
// 	aggregate, err := h.Tracks.GetNormalizedJourneyAggregate(ownerID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		log.Println("Error collecting aggregate: ", err)
// 		return
// 	}

// 	// Write gif back to client
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(aggregate)
// }

// ~=~=~=~=~=~=~=~=
// Kpis
// ~=~=~=~=~=~=~=~=

func (h *Handler) newKpi(w http.ResponseWriter, r *http.Request) {
	var kpi app.Kpi

	// Parse body
	err := json.NewDecoder(r.Body).Decode(&kpi)
	if err != nil {
		http.Error(w, invalidRequestError, http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Store KPI
	newKpiID, err := h.Kpis.Store(kpi)
	if err != nil {
		http.Error(w, internalError, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Response
	s := strconv.FormatInt(newKpiID, 10)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, s)
}

func (h *Handler) deleteKpi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]
	ownerIDInterface := r.Context().Value(ContextKeyOwnerID)

	ownerID, ok := ownerIDInterface.(int64)
	if !ok {
		http.Error(w, "owner id error", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, "id error", http.StatusBadRequest)
		return
	}

	// Delete KPI
	deleted, err := h.Kpis.Delete(id, ownerID)
	if err != nil {
		http.Error(w, internalError, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Response
	s := strconv.FormatInt(deleted, 10)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, s)
}

func (h *Handler) listKpis(w http.ResponseWriter, r *http.Request) {
	ownerID := r.Context().Value(ContextKeyOwnerID).(int64)

	// Get Kpis
	kpis, err := h.Kpis.FindByOwnerID(ownerID)
	if err != nil {
		http.Error(w, internalError, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Get aggregates for the kpi
	for i, kpi := range kpis {
		// Get aggregate data
		aggregate, err := h.Tracks.GetNormalizedJourneyAggregate(kpi.OwnerID, "campaign_name", kpi.PatternMatchColumnName, kpi.PatternMatchRowValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error collecting aggregate: ", err)
			return
		}
		kpis[i].CampaignNameJourneyAggregate = aggregate
	}

	// Format
	if kpis == nil {
		kpis = []app.Kpi{}
	}

	// Response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kpis)
}
