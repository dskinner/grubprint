package api

import (
	"fmt"
	"net/http"
	"time"

	"dasa.cc/food/httputil"
	"dasa.cc/food/keystore"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
)

// defaultTimeout is the duration before a request is cancelled.
var defaultTimeout = 1 * time.Second

// handler implements http.Handler and implementations must map mux.Vars(r)
// to the appropriate datastore service. Handlers run in a separate goroutine
// that will timeout by the duration of defaultTimeout.
type handler func(ctx context.Context, r *http.Request) (interface{}, error)

func (h handler) ServeHTTP(resp http.ResponseWriter, r *http.Request) {
	w := &httputil.ResponseWriter{resp, http.StatusOK}

	tr := trace.New("api."+mux.CurrentRoute(r).GetName(), r.URL.Path)
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

	if err := keystore.VerifyRequest(r); err != nil {
		switch err := err.(type) {
		case keystore.Error:
			keystore.HandleError(w, err)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		tr.LazyPrintf("AUTH %v\n", err)
		tr.SetError()
		return
	}

	if err := serveContent(ctx, w, r, h); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		tr.LazyPrintf("ERROR %v\n", err)
		tr.SetError()
	}
}

func serveContent(ctx context.Context, w http.ResponseWriter, r *http.Request, h handler) error {
	if httputil.CheckModified(ctx, w, r) {
		return nil
	}

	out := make(chan interface{})
	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				out <- fmt.Errorf("%v", rec)
			}
		}()
		v, err := h(ctx, r)
		if err != nil {
			out <- err
		} else {
			out <- v
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case v := <-out:
		switch t := v.(type) {
		case error:
			return t
		default:
			return httputil.WriteJSON(w, t)
		}
	}
}
