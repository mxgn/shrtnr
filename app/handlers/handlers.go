package handlers

import (
	"log"
	"net/http"

	"github.com/mxgn/url-shrtnr/app/storages"
)

func EncodeHandler(storage storages.IStorage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {

		if url := r.PostFormValue("url"); url != "" {

			shortResult := storage.Save(url)
			linkText := "http://localhost/" + shortResult
			response := "<a href=\"" + linkText + "\">" + linkText + "</a>"

			w.Write([]byte(response))
		} else {
			w.Write([]byte("Error: no URL provided/"))
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

func RedirectHandler(storage storages.IStorage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Path[len("/"):]

		if code == "" {
			http.Redirect(w, r, "/add/", 303)
		}

		url, err := storage.Load(code)
		if err != nil {
			log.Println("")
		}
		log.Println("Long url from database:", url)

		http.Redirect(w, r, "http://"+string(url), http.StatusSeeOther)
	}
	return http.HandlerFunc(handleFunc)
}

func Handler404(storage storages.IStorage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	return http.HandlerFunc(handleFunc)
}
