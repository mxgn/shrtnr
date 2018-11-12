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

	log.SetFlags(log.Lshortfile &^ (log.Ldate | log.Ltime))

	db := &storages.DbIface{}
	db.Init()
	db.CreateSchema()

	r := mux.NewRouter()

	r.Handle("/add/", http.StripPrefix("/add/", http.FileServer(http.Dir("/var/www"))))
	r.HandleFunc("/", handlers.UrlRedirect)
	r.HandleFunc("/add", handlers.UrlAdd).Methods(http.MethodPost)
	r.HandleFunc("/{^[A-Za-z0-9]+$}", handlers.UrlRedirect).Methods(http.MethodGet)

	http.Handle("/", r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
