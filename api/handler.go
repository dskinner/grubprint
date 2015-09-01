package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
)

var (
	// modtime is used for Last-Modified header since usda data is static.
	modtime = time.Now().UTC()

	// defaultTimeout is the duration before a request is cancelled.
	defaultTimeout = 1 * time.Second
)

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

// handler implements http.Handler and implementations must map mux.Vars(r)
// to the appropriate datastore service. Handlers run in a separate goroutine
// that will timeout by the duration of defaultTimeout.
type handler func(ctx context.Context, r *http.Request) (interface{}, error)

func (h handler) ServeHTTP(resp http.ResponseWriter, r *http.Request) {
	w := &responseWriter{resp, http.StatusOK}

	tr := trace.New("api."+mux.CurrentRoute(r).GetName(), r.URL.Path)
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

	if err := serveContent(ctx, w, r, h); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		tr.LazyPrintf("ERROR %v\n", err)
		tr.SetError()
	}
}

// checkModified inspects If-Modified-Since and returns whether request is now complete.
func checkModified(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {
	tr, _ := trace.FromContext(ctx)

	modsince := r.Header.Get("If-Modified-Since")
	if modsince == "" {
		return false
	}

	t, err := time.Parse(http.TimeFormat, modsince)
	if err != nil {
		tr.LazyPrintf("ERROR If-Modified-Since: %v\n", err)
		return false
	}

	if modtime.Before(t.Add(1 * time.Second)) {
		w.WriteHeader(http.StatusNotModified)
		return true
	}

	return false
}

func serveContent(ctx context.Context, w http.ResponseWriter, r *http.Request, h handler) error {
	if checkModified(ctx, w, r) {
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
			return write(w, t)
		}
	}
}

func write(w http.ResponseWriter, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age:86400")
	w.Header().Set("Last-Modified", modtime.Format(http.TimeFormat))

	_, err = w.Write(data)
	return err
}
