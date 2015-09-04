package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"grubprint.io/datastore"
	"grubprint.io/router"
)

var store = datastore.New()

func Handler() http.Handler {
	r := router.New()
	r.Get(router.Foods).Handler(handler(foods))
	r.Get(router.Weights).Handler(handler(weights))
	r.Get(router.Nutrients).Handler(handler(nutrients))
	return r
}

func foods(ctx context.Context, r *http.Request) (interface{}, error) {
	return store.Foods.Search(mux.Vars(r)["q"])
}

func weights(ctx context.Context, r *http.Request) (interface{}, error) {
	return store.Weights.ByFoodId(mux.Vars(r)["id"])
}

func nutrients(ctx context.Context, r *http.Request) (interface{}, error) {
	return store.Nutrients.ByFoodId(mux.Vars(r)["id"])
}
