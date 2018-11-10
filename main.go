package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mxgn/url-shrtnr/app/handlers"
	"github.com/mxgn/url-shrtnr/app/storages"
)

func main() {
	log.SetFlags(log.Lshortfile &^ (log.Ldate | log.Ltime))

	storage := &storages.Pgdb{}
	storage.Init()
	storage.CreateSchema()

	fs := http.FileServer(http.Dir("/var/www"))
	http.Handle("/add/", http.StripPrefix("/add/", fs))

	http.Handle("/", handlers.RedirectHandler(storage))
	http.Handle("/add", handlers.EncodeHandler(storage))
	http.Handle("/favicon.ico", handlers.Handler404(storage))

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
