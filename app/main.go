package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mxgn/url-shrtnr/app/handlers"
	"github.com/mxgn/url-shrtnr/app/models"
	"github.com/mxgn/url-shrtnr/app/storages"
)

func main() {

	log.SetFlags(log.LstdFlags &^ (log.Ldate | log.Ltime))
	// log.SetFlags(log.Lshortfile)

	app := &models.AppConfig{Debug: true}
	app.Init()

	storages.Pgdb = storages.Init(app)
	// storages.Pgdb.CreateSchema()

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir(app.StaticDir))
	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", fs)).
		Methods("GET")

	r.HandleFunc("/", handlers.UrlRedirect).
		Methods(http.MethodGet)

	r.HandleFunc("/api/add", handlers.UrlAdd).
		Methods(http.MethodPost)

	r.HandleFunc("/{^[A-Za-z0-9]+$}", handlers.UrlRedirect).
		Methods("GET")

	http.Handle("/", r)
	if err := http.ListenAndServe(":"+app.Port, nil); err != nil {
		log.Fatal(err)
	}

}
