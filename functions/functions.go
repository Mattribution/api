package functions

import (
	"net/http"

	"github.com/mattribution/api/common"
)

func NewTrack(w http.ResponseWriter, r *http.Request) {
	msg := common.GenerateMessage()
	// our code will go here
	w.WriteHeader(200)
	w.Write([]byte(msg))
}
