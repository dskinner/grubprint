package router

import "github.com/gorilla/mux"

const (
	Foods     = "Foods"
	Weights   = "Weights"
	Nutrients = "Nutrients"
)

func New() *mux.Router {
	r := mux.NewRouter()

	r.Path("/foods/{q}").Methods("GET").Name(Foods)
	r.Path("/weights/{id}").Methods("GET").Name(Weights)
	r.Path("/nutrients/{id}").Methods("GET").Name(Nutrients)

	return r
}
