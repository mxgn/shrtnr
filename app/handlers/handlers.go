// Package handlers provides HTTP request handlers.
package handlers

import (
	"log"
	"net/http"

	"github.com/mxgn/url-shrtnr/app/storages"
)

//EncodeHandler (storage storages.IStorage) http.Handler {
func EncodeHandler(storage storages.IStorage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		if url := r.PostFormValue("url"); url != "" {
			log.Println("URL IS: %d " + url)
			w.Write([]byte(storage.Save(url)))
		}
	}
	return http.HandlerFunc(handleFunc)
}

//DecodeHandler (storage storages.IStorage) http.Handler {
func DecodeHandler(storage storages.IStorage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Path[len("/dec/"):]

		url, err := storage.Load(code)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("URL Not Found. Error: " + err.Error() + "\n"))
			return
		}

		w.Write([]byte(url))
	}

	return http.HandlerFunc(handleFunc)
}

//RedirectHandler (storage storages.IStorage) http.Handler {
func RedirectHandler(storage storages.IStorage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Path[len("/"):]

		url, err := storage.Load(code)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("URL Not Found. Error: " + err.Error() + "\n"))
			return
		}

		http.Redirect(w, r, string(url), 301)
	}

	return http.HandlerFunc(handleFunc)
}
