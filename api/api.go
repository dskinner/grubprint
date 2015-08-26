package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"dasa.cc/food/datastore"
	"dasa.cc/food/router"
	"github.com/gorilla/mux"
)

var store = datastore.New()

func Router() *mux.Router {
	r := router.New()
	r.Get(router.Foods).Handler(handler(foods))
	r.Get(router.Weights).Handler(handler(weights))
	r.Get(router.Nutrients).Handler(handler(nutrients))
	return r
}

type handler func(http.ResponseWriter, *http.Request) error

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %s", err)
		log.Println(err)
	}
}

func write(w http.ResponseWriter, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	_, err = w.Write(data)
	return err
}

func foods(w http.ResponseWriter, r *http.Request) error {
	q := mux.Vars(r)["q"]
	m, err := store.Foods.Search(q)
	if err != nil {
		return err
	}
	return write(w, m)
}

func weights(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	m, err := store.Weights.ByFoodId(id)
	if err != nil {
		return err
	}
	return write(w, m)
}

func nutrients(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	m, err := store.Nutrients.ByFoodId(id)
	if err != nil {
		return err
	}
	return write(w, m)
}
