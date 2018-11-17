package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mxgn/url-shrtnr/app/storages"
)

func main() {

	log.SetFlags(log.LstdFlags &^ (log.Ldate | log.Ltime))
	log.SetFlags(log.Lshortfile)

	app := &AppConfig{Debug: true}

	app.Init()

	app.Storage = storages.Init(true)

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir(app.StaticDir))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs)).Methods("GET")

	r.HandleFunc("/", UrlRedirect(app)).Methods("GET")
	r.HandleFunc("/api/add", UrlAdd(app)).Methods("POST")
	r.HandleFunc("/{^[A-Za-z0-9]+$}", UrlRedirect(app)).Methods("GET")

	http.Handle("/", r)
	if err := http.ListenAndServe(":"+app.Port, nil); err != nil {
		log.Fatal(err)
	}

}
