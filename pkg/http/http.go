package http

import (
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mattribution/api/pkg/api"
)

var gif = []byte{
	71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 0, 0, 0,
	255, 255, 255, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0,
	1, 0, 1, 0, 0, 2, 1, 68, 0, 59,
}

type Handler struct {
	TrackService api.TrackService
}

func NewHandler(trackService api.TrackService) Handler {
	return Handler{
		TrackService: trackService,
	}
}

func (h *Handler) NewTrack(w http.ResponseWriter, r *http.Request) {
	// Get pixel data from client
	v := r.URL.Query()
	rawEvent := v.Get("data")
	data, err := base64.StdEncoding.DecodeString(rawEvent)
	if err != nil {
		panic(err)
	}

	// Unmarshal
	track := api.Track{}
	if err := json.Unmarshal(data, &track); err != nil {
		panic(err)
	}

	// Grab IP
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	track.IP = ip

	// Store
	_, err = h.TrackService.StoreTrack(track)
	if err != nil {
		panic(err)
	}

	// Write gif back to client
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Content-Type", "image/gif")
	w.Write(gif)
}

// Serve http
func (h *Handler) Serve(addr string) error {
	// Setup mux
	r := mux.NewRouter()
	r.HandleFunc("/v1/pixel/track", h.NewTrack).Methods("GET")
	return http.ListenAndServe(addr, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r))
}
