package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mxgn/url-shrtnr/app/config"
	"github.com/mxgn/url-shrtnr/app/handlers"
	"github.com/mxgn/url-shrtnr/app/storage/postgre"
)

func main() {

	log.SetFlags(log.LstdFlags &^ (log.Ldate | log.Ltime))
	log.SetFlags(log.Lshortfile)

	app := &config.AppContext{Debug: false}

	app.ReadConfig()

	app.DB = postgre.Init(app)

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir(app.StaticDir))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs)).Methods("GET")

	r.HandleFunc("/", handlers.UrlRedirect(app)).Methods("GET")
	r.HandleFunc("/api/add", handlers.UrlAdd(app)).Methods("POST")
	r.HandleFunc("/{^[A-Za-z0-9]+$}", handlers.UrlRedirect(app)).Methods("GET")

	http.Handle("/", r)
	if err := http.ListenAndServe(":"+app.Port, nil); err != nil {
		log.Fatal(err)
	}

}
