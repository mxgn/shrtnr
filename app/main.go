package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mxgn/url-shrtnr/app/handlers"
	"github.com/mxgn/url-shrtnr/app/storages"
)

func main() {

	storage := &storages.Redis{}
	if err := storage.Init(); err != nil {
		log.Fatal(err)
	}

	http.Handle("/", handlers.RedirectHandler(storage))
	http.Handle("/enc/", handlers.EncodeHandler(storage))
	http.Handle("/dec/", handlers.DecodeHandler(storage))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
