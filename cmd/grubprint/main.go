package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	_ "expvar"
	_ "net/http/pprof"

	"grubprint.io/api"
	"grubprint.io/app"
	"grubprint.io/keystore"
)

var (
	flagAddr   = flag.String("addr", ":8080", "address to listen on")
	flagStatic = flag.String("static", "app/polymer/dist", "directory of static resources")
	flagKeygen = flag.Bool("keygen", false, "generate new key pair, write to disk, and return")
)

func main() {
	flag.Parse()

	if *flagKeygen {
		pub, priv, err := keystore.Keygen()
		if err != nil {
			log.Fatal(err)
		}
		if err := ioutil.WriteFile("id_rsa", priv, 0644); err != nil {
			log.Fatal(err)
		}
		if err := ioutil.WriteFile("id_rsa.pub", pub, 0644); err != nil {
			log.Fatal(err)
		}
		log.Println("new key pair generated")
		return
	}

	bin, err := ioutil.ReadFile("id_rsa.pub")
	if err != nil {
		log.Fatal(err)
	}
	if err := keystore.Set("oauth2@keystore", bin); err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(*flagStatic); os.IsNotExist(err) {
		wd, _ := os.Getwd()
		log.Fatalf("static directory missing: %s\n", path.Join(wd, *flagStatic))
	}

	m := http.NewServeMux()
	m.Handle("/debug/", http.DefaultServeMux)
	m.Handle("/oauth2/token", keystore.TokenHandler)
	m.Handle("/api/", http.StripPrefix("/api", api.Handler()))
	m.Handle("/app/", http.StripPrefix("/app", app.Handler()))
	// m.Handle("/dist/", http.StripPrefix("/dist", http.FileServer(http.Dir(*flagStatic))))
	m.Handle("/", http.FileServer(http.Dir(*flagStatic)))

	log.Println("listening on", *flagAddr)
	log.Fatal("ListenAndServe:", http.ListenAndServe(*flagAddr, m))
}
