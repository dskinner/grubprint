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
)

var cl = client.New(nil)

var defaultTimeout = 1 * time.Second

func Handler() http.Handler {
	r := router.New()
	r.Get(router.Index).Handler(handler(index))
	r.Get(router.Foods).Handler(handler(foods))
	return r
}

func index(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return httputil.WriteHTML(w, "index.html", "Hello World")
}

func foods(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	m, err := cl.Foods.Search(mux.Vars(r)["q"])
	if err != nil {
		return fmt.Errorf("Client: %v", err)
	}
	return httputil.WriteHTML(w, "foods.html", m)
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
