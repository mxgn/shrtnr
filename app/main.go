package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mxgn/url-shrtnr/app/handlers"
	"github.com/mxgn/url-shrtnr/app/storages"
)

var db *sql.DB

func main() {

	log.SetFlags(log.LstdFlags &^ (log.Ldate | log.Ltime))
	// log.SetFlags(log.Lshortfile)

	db := &storages.DbIface{}
	db.Init()
	db.CreateSchema()

	cfg := Init(true)

	log.Println(`cfg.debug=`, cfg.debug)

	handlers.AppConfig := cfg

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir(cfg.staticDir))
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
