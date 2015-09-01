package main

import (
	"flag"
	"log"
	"net/http"

	_ "expvar"
	_ "net/http/pprof"

	"dasa.cc/food/api"
	"dasa.cc/food/app"
	"dasa.cc/food/datastore"
)

var (
	flagAddr = flag.String("addr", ":8080", "address to listen on")
)

func main() {
	flag.Parse()

	datastore.Connect()

	m := http.NewServeMux()
	m.Handle("/debug/", http.DefaultServeMux)
	m.Handle("/api/", http.StripPrefix("/api", api.Handler()))
	m.Handle("/", app.Handler())

	log.Println("listening on", *flagAddr)
	log.Fatal("ListenAndServe:", http.ListenAndServe(*flagAddr, m))
}
