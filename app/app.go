package app

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"dasa.cc/food/client"
	"dasa.cc/food/router"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
)

var cl = client.New(nil)

var defaultTimeout = 1 * time.Second

func Handler() http.Handler {
	r := router.New()

	// name the path for tracing
	r.Path("/").Methods("GET").Name("Index")
	r.Get("Index").Handler(handler(index))

	r.Get(router.Foods).Handler(handler(foods))

	return r
}

// responseWriter wraps http.ResponseWriter to provide status for logging.
//
// Calling ResponseWriter.Write does not actually call WriteHeader until
// response is flushed. Suggested use is to initialize status with
// http.StatusOK and only reference this value once the response is
// complete.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

type handler func(context.Context, http.ResponseWriter, *http.Request) error

func (h handler) ServeHTTP(resp http.ResponseWriter, r *http.Request) {
	w := &responseWriter{resp, http.StatusOK}

	tr := trace.New("app."+mux.CurrentRoute(r).GetName(), r.URL.Path)
	defer tr.Finish()

	tr.LazyPrintf("HTTP %s %s\n", r.Method, r.URL.Path)
	defer func() { tr.LazyPrintf("END HTTP %v %s\n", w.status, http.StatusText(w.status)) }()

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

	if err := h(ctx, w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		tr.LazyPrintf("ERROR %v\n", err)
		tr.SetError()
	}
}

func index(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	t, err := template.ParseFiles("app/templates/index.html")
	if err != nil {
		return err
	}
	return t.Execute(w, nil)
}

func foods(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	m, err := cl.Foods.Search(mux.Vars(r)["q"])
	if err != nil {
		return fmt.Errorf("Client: %v", err)
	}
	return write(w, m)
}

func write(w http.ResponseWriter, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	_, err = w.Write(data)
	return err
}
