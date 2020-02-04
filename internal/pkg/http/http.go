package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mattribution/api/internal/pkg/app"
	"github.com/mattribution/api/internal/pkg/postgres"
)

// Handler combines services to create an http handler
type Handler struct {
	Tracks app.Tracks
}

func (h Handler) NewTrack(w http.ResponseWriter, req *http.Request) {
	h.Tracks.Store(app.Track{})

	w.WriteHeader(200)
	fmt.Fprint(w, "Track stored!")
}

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r := mux.NewRouter()
	s := r.PathPrefix("/v1").Subrouter()
	s.HandleFunc("/track", h.NewTrack).Methods(http.MethodGet)
	s.ServeHTTP(w, req)
}

// ListenAndServe listens and serves
func (h Handler) ListenAndServe() {
	http.ListenAndServe(":3000", h)
}

// CloudFunctionEntryPoint represents cloud function entry point
func CloudFunctionEntryPoint(w http.ResponseWriter, r *http.Request) {
	tracks := postgres.Tracks{}

	handler := Handler{
		Tracks: &tracks,
	}

	handler.ServeHTTP(w, r)
}
