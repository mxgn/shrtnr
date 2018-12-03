package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mxgn/url-shrtnr/app/config"
	"github.com/mxgn/url-shrtnr/app/handlers"
	"github.com/mxgn/url-shrtnr/app/helpers"
	"github.com/mxgn/url-shrtnr/app/storage/postgre"

	log "github.com/mxgn/seelog"
)

func main() {

	// log.SetFlags(log.LstdFlags &^ (log.Ldate | log.Ltime))
	// log.SetFlags(log.Lshortfile)

	testConfig := `
	<seelog minlevel="trace" type="sync">
		<outputs formatid="main">
			<console/>
		</outputs>
		<formats>
			<format id="main" format="%File: %FuncShort %Line: %Msg%n"/>
		</formats>
	</seelog>`
	logger, _ := log.LoggerFromConfigAsBytes([]byte(testConfig))
	log.ReplaceLogger(logger)

	c := &config.AppContext{Log: logger}
	c.Init()

	handlers.Init(c)
	helpers.Init(c)

	c.DB = postgre.Init(c)
	// postgre.CreateSchema()

	r := mux.NewRouter()
	fs := http.FileServer(http.Dir(c.StaticDir))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs)).Methods("GET")
	r.HandleFunc("/", handlers.UrlRedirectIndex(c)).Methods("GET")
	r.HandleFunc("/api/add/", handlers.UrlAdd(c)).Methods("POST")
	r.HandleFunc("/{^[A-Za-z0-9]+$}", handlers.UrlRedirect(c)).Methods("GET")
	http.Handle("/", r)
	log.Critical(http.ListenAndServe(":"+c.Port, nil))

	// r := mux.NewRouter()
	// fs := http.FileServer(http.Dir(c.StaticDir))
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs)).Methods("GET")
	// http.Handle("/static/", r)
	// http.Handle("/api/add/", handlers.Handler{Ctx: c, H: handlers.AddNew})
	// http.Handle("/", handlers.Handler{Ctx: c, H: handlers.RedirectNew})
	// log.Fatal(http.ListenAndServe(":"+c.Port, nil))

}
