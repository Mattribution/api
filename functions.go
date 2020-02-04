package functions

import (
	"log"
	"net/http"

	"github.com/mattribution/api/common"
	"github.com/mattribution/api/internal/app"
)

// func exampleMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Our middleware logic goes here...
// 		next.ServeHTTP(w, r)
// 	})
// }

func NewTrack(w http.ResponseWriter, r *http.Request) {
	t := app.Track{}
	log.Println(t)

	msg := common.GenerateMessage()
	// our code will go here
	w.WriteHeader(200)
	w.Write([]byte(msg))
}

// var (
// 	newTrackHandler    = http.HandlerFunc(newTrack)
// 	NewTrackAuthorized = exampleMiddleware(newTrackHandler)
// )
