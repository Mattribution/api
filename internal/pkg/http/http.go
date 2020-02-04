package http

import (
	"fmt"
	"net/http"

	"github.com/mattribution/api/internal/pkg/app"
)

// Handler combines services to create an http handler
type Handler struct {
	tracks app.Tracks
}

func (h Handler) NewTrack(w http.ResponseWriter, r *http.Request) {
	h.tracks.Store(app.Track{})

	w.WriteHeader(200)
	fmt.Fprint(w, "Track stored!")
}
