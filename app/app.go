package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"grubprint.io/client"
	"grubprint.io/httputil"
	"grubprint.io/router"
	"grubprint.io/usda"
)

var cl *client.Client

var defaultTimeout = 1 * time.Second

func Handler() http.Handler {
	// TODO work out initialization of client regarding keys
	if cl == nil {
		cl = client.New(nil)
	}
	r := router.New()
	r.Get(router.Index).Handler(handler(index))
	r.Get(router.Foods).Handler(handler(foods))
	r.Get(router.Nutrients).Handler(handler(nutrients))
	return r
}

func index(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return httputil.WriteHTML(w, "index.html", nil)
}

func foods(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	foods, err := cl.Foods.Search(mux.Vars(r)["q"])
	if err != nil {
		return fmt.Errorf("Client: %v", err)
	}
	return httputil.WriteHTML(w, "foods.html", foods)
}

func nutrients(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	type item struct {
		Food      *usda.Food
		Weights   []*usda.Weight
		Nutrients *usda.Nutrients
	}

	food, err := cl.Foods.ById(mux.Vars(r)["id"])
	if err != nil {
		return fmt.Errorf("Client: %v", err)
	}
	weights, err := cl.Weights.ByFoodId(food.Id)
	if err != nil {
		return fmt.Errorf("Client: %v", err)
	}
	nutrients, err := cl.Nutrients.ByFoodId(food.Id)
	if err != nil {
		return fmt.Errorf("Client: %v", err)
	}

	return httputil.WriteHTML(w, "detail.html", item{food, weights, usda.NewNutrients(nutrients...)})
}

type handler func(context.Context, http.ResponseWriter, *http.Request) error

func (h handler) ServeHTTP(resp http.ResponseWriter, r *http.Request) {
	w := &httputil.ResponseWriter{resp, http.StatusOK}

	tr := trace.New("app."+mux.CurrentRoute(r).GetName(), r.URL.Path)
	defer tr.Finish()

	tr.LazyPrintf("HTTP %s %s\n", r.Method, r.URL.Path)
	defer func() { tr.LazyPrintf("END HTTP %v %s\n", w.Status, http.StatusText(w.Status)) }()

	defer func() {
		if rec := recover(); rec != nil {
			http.Error(w, fmt.Sprintf("%v", rec), http.StatusInternalServerError)
			tr.LazyPrintf("PANIC %v\n", rec)
			tr.SetError()
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	ctx = trace.NewContext(ctx, tr)

	if httputil.CheckModified(ctx, w, r) {
		return
	}

	if err := h(ctx, w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		tr.LazyPrintf("ERROR %v\n", err)
		tr.SetError()
	}
}
