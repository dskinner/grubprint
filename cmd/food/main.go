package main

import (
	"log"
	"net/http"

	"dasa.cc/food/api"
	"dasa.cc/food/datastore"
)

func main() {
	datastore.Connect()
	log.Fatal(http.ListenAndServe(":8080", api.Router()))
}
