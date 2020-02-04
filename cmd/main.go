package main

import (
	"github.com/mattribution/api/internal/pkg/http"
	"github.com/mattribution/api/internal/pkg/postgres"
)

func main() {
	tracks := postgres.Tracks{}

	handler := http.Handler{
		Tracks: &tracks,
	}

	handler.ListenAndServe()
}
