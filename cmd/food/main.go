package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"

	_ "expvar"
	_ "net/http/pprof"

	"dasa.cc/food/api"
	"dasa.cc/food/app"
	"dasa.cc/food/datastore"
)

var (
	flagAddr   = flag.String("addr", ":8080", "address to listen on")
	flagStatic = flag.String("static", "app/static", "directory of static resources")
)

func main() {
	flag.Parse()

	if _, err := os.Stat(*flagStatic); os.IsNotExist(err) {
		wd, _ := os.Getwd()
		log.Fatalf("static directory missing: %s\n", path.Join(wd, *flagStatic))
	}

	datastore.Connect()

	m := http.NewServeMux()
	m.Handle("/debug/", http.DefaultServeMux)
	m.Handle("/api/", http.StripPrefix("/api", api.Handler()))
	m.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(*flagStatic))))
	m.Handle("/", app.Handler())

	log.Println("listening on", *flagAddr)
	log.Fatal("ListenAndServe:", http.ListenAndServe(*flagAddr, m))
}
