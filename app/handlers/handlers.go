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

		// reqbyte, _ := httputil.DumpRequest(r, true)
		// log.Println(string(reqbyte))

		// r.ParseForm() //анализ аргументов,
		// fmt.Println("ENC REQUEST RORM:", r.Form)

		if url := r.PostFormValue("url"); url != "" {

			shortResult := storage.Save(url)

			response := "http://localhost/" + shortResult

			w.Write([]byte(response))
		} else {
			w.Write([]byte("TETREWTWS!!"))
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

		// reqbyte, _ := httputil.DumpRequest(r, true)
		// log.Println("!Request Byte:\n\n<--", string(reqbyte), "<---")

		// r.ParseForm() //анализ аргументов,
		// fmt.Println("URL FORM:", r.Form)

		code := r.URL.Path[len("/"):]

		log.Println("Short code requested, req.URL.Path:", code)

		if code == "" {
			http.Redirect(w, r, "/notfound", 303)
		}

		url, err := storage.Load(code)
		if err != nil {
			// do nothing:

			// w.WriteHeader(http.StatusNotFound)
			// w.Write([]byte("URL Not Found. Error: " + err.Error() + "\n"))
			// return
		}
		log.Println("Long url from database:", url)

		http.Redirect(w, r, "http://"+string(url), http.StatusSeeOther)
	}

	return http.HandlerFunc(handleFunc)
}

//RedirectHandler (storage storages.IStorage) http.Handler {
func NullHandler(storage storages.IStorage) http.Handler {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	return http.HandlerFunc(handleFunc)
}
