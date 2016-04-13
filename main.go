package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	_ "expvar"
	_ "net/http/pprof"

	"grubprint.io/api"
	"grubprint.io/app"
	"grubprint.io/client"
	"grubprint.io/datastore"
	"grubprint.io/httputil"
	"grubprint.io/keystore"
)

var (
	flagAddr   = flag.String("addr", ":8080", "address to listen on")
	flagAssets = flag.String("assets", "assets", "directory of static resources")
	flagDB     = flag.String("db", "usda.db", "location of db file")
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

	if _, err := os.Stat(*flagAssets); os.IsNotExist(err) {
		wd, _ := os.Getwd()
		log.Fatalf("assets directory missing: %s\n", path.Join(wd, *flagAssets))
	}
	httputil.AssetsDir = *flagAssets
	client.AssetsDir = *flagAssets

	bin, err := ioutil.ReadFile(filepath.Join(*flagAssets, "id_rsa.pub"))
	if err != nil {
		log.Fatal(err)
	}
	if err := keystore.Set("oauth2@keystore", bin); err != nil {
		log.Fatal(err)
	}

	datastore.Connect(*flagDB)

	m := http.NewServeMux()
	m.Handle("/debug/", http.DefaultServeMux)
	m.Handle("/oauth2/token", keystore.TokenHandler)
	m.Handle("/api/", http.StripPrefix("/api", api.Handler()))
	m.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir(filepath.Join(*flagAssets, "public")))))
	m.Handle("/", app.Handler())

	log.Println("listening on", *flagAddr)
	log.Fatal("ListenAndServe:", http.ListenAndServe(*flagAddr, m))
}
