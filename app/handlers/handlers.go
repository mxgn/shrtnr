package handlers

import (
	"log"
	"net/http"

	"github.com/mxgn/url-shrtnr/app/storages"
)

func UrlRedirect(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Path[len("/"):]
	log.Println("r.URL.Path[len(\"/\"):]   :", code)

	if code == "" {
		http.Redirect(w, r, "/add/", 303)
		return
	}

	url, err := storages.Pgdb.Load(code)
	if err != nil {
		log.Println("")
	}
	log.Println("Long url from database:", url)
	http.Redirect(w, r, "http://"+string(url), http.StatusSeeOther)
}

func UrlAdd(w http.ResponseWriter, r *http.Request) {

	test := r.URL.Query().Get("url")
	log.Println("r.URL.QueryGet(\"url\"):", test)

	var response string

	if url := r.PostFormValue("url"); url != "" {

		shortResult := storages.Pgdb.Save(url)
		linkUrl := "http://localhost/" + shortResult
		response = "<a href=\"" + linkUrl + "\">" + linkUrl + "</a>"

	} else {
		response = "Url add Handler - no form parameter url provided"
	}

	w.Write([]byte(response))
}
