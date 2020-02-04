package postgres

import (
	"log"

	"github.com/mattribution/api/internal/pkg/app"
)

type Tracks struct{}

// Store stores a given track in the database
func (t *Tracks) Store(track app.Track) (int64, error) {
	log.Println("TODO: Store")
	return 1, nil
}
