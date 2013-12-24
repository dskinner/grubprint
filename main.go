package main

import (
	"dasa.cc/dae"
	"dasa.cc/dae/handler"
	"dasa.cc/food/usda"
	"flag"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

var (
	dbinit = flag.Bool("dbinit", false, "init the database and exit")
	router = mux.NewRouter()
)

func init() {
	flag.Parse()
}

func index(w http.ResponseWriter, r *http.Request) *handler.Error {
	t, err := template.ParseFiles("views/index.html")
	if err != nil {
		return handler.NewError(err, 500, "Failed to parse index.html")
	}
	if err := t.Execute(w, nil); err != nil {
		return handler.NewError(err, 500, "Failed to execute template")
	}
	return nil
}

func main() {
	if *dbinit {
		usda.LoadAll()
		return
	}
	dae.RegisterFileServer("res/")
	// dae.ServeFile("/favicon.ico", "res/favicon.ico")

	router.Handle("/", handler.New(index))

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe("localhost:8090", nil))
}
