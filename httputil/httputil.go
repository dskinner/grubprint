package httputil

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"
)

// modtime is used for Last-Modified header since usda data is static.
var modtime = time.Now().UTC()

// CheckModified inspects If-Modified-Since and returns whether request is now complete.
func CheckModified(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {
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

// ResponseWriter wraps http.ResponseWriter to provide status for logging.
//
// Calling ResponseWriter.Write does not actually call WriteHeader until
// response is flushed. Suggested use is to initialize status with
// http.StatusOK and only reference this value once the response is
// complete.
type ResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (w *ResponseWriter) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

func WriteJSON(w http.ResponseWriter, v interface{}) error {
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

func WriteHTML(w http.ResponseWriter, name string, v interface{}) error {
	tmpl, err := template.New(name).ParseFiles(filepath.Join("app", "assets", "templates", name))
	if err != nil {
		return err
	}
	return tmpl.Execute(w, v)
}
