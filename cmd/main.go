package main

import (
	functions "github.com/mattribution/api"
	"net/http"
)

func main() {
	http.HandleFunc("/", functions.FunctionsEntrypoint)
	panic(http.ListenAndServe(":3001", nil))
}
