package router

import "github.com/gorilla/mux"

const (
	Index     = "Index"
	Food      = "Food"
	Foods     = "Foods"
	Weights   = "Weights"
	Nutrients = "Nutrients"
)

func New() *mux.Router {
	r := mux.NewRouter()

	r.Path("/").Methods("GET").Name(Index)
	r.Path("/foods/{q}").Methods("GET").Name(Foods)
	r.Path("/food/{id}").Methods("GET").Name(Food)
	r.Path("/weights/{id}").Methods("GET").Name(Weights)
	r.Path("/nutrients/{id}").Methods("GET").Name(Nutrients)

	return r
}
