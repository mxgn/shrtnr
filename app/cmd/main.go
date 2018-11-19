package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mxgn/url-shrtnr/app/config"
	"github.com/mxgn/url-shrtnr/app/handlers"
	"github.com/mxgn/url-shrtnr/app/storage/postgre"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {

	// log.SetFlags(log.LstdFlags &^ (log.Ldate | log.Ltime))
	// log.SetFlags(log.Lshortfile)

	log.Out = os.Stdout
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors:    false,
		FullTimestamp:    false,
		DisableTimestamp: true,
		QuoteEmptyFields: true,
	})

	app := &config.AppContext{Debug: true, Log: log}

	app.ReadConfig()

	app.DB = postgre.Init(app)
	// postgre.CreateSchema()

	r := mux.NewRouter()
	fs := http.FileServer(http.Dir(app.StaticDir))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs)).Methods("GET")
	r.HandleFunc("/", handlers.UrlRedirectIndex(app)).Methods("GET")
	r.HandleFunc("/api/add/", handlers.UrlAdd(app)).Methods("POST")
	r.HandleFunc("/{^[A-Za-z0-9]+$}", handlers.UrlRedirect(app)).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":"+app.Port, nil))

	// r := mux.NewRouter()
	// fs := http.FileServer(http.Dir(app.StaticDir))
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs)).Methods("GET")
	// http.Handle("/static/", r)
	// http.Handle("/api/add/", handlers.Handler{Ctx: app, H: handlers.AddNew})
	// http.Handle("/", handlers.Handler{Ctx: app, H: handlers.RedirectNew})
	// log.Fatal(http.ListenAndServe(":"+app.Port, nil))

}
