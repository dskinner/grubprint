package api

import (
	"net/http"

	"dasa.cc/food/datastore"
	"dasa.cc/food/router"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

var store = datastore.New()

func Router() http.Handler {
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
