package handlers

import (
	"log"
	"net/http"

	"github.com/mxgn/url-shrtnr/app/config"
)

var err error

func UrlRedirect(c *config.AppContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if c.Debug {
			log.Printf("GOT NEW REQUEST: Req.Host=%s, Req.URL=%s\n", r.Host, r.URL)
		}

		code := r.URL.Path[len("/"):]

		url, err := c.DB.GetLongUrl(code)
		if err != nil {
			log.Println("UrlRedirect, short code not found:", code)
			http.Redirect(w, r, "URL "+code+" NotFound", http.StatusNotFound)
			return
		}

		log.Println("Long url from database:", url)
		http.Redirect(w, r, "http://"+string(url), http.StatusSeeOther)
	}
}

func UrlAdd(app *config.AppContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println(`r.PostFormValue("url"):`, r.PostFormValue("url"))

		var response string

		if url := r.PostFormValue("url"); url != "" {

			shortUrl := ""
			shortUrl, err = app.DB.AddLongUrl(url)
			if err != nil {
				log.Panicln("CANNOT ADD URL:", err)
			}

			linkUrl := "http://localhost/" + shortUrl
			response = "<a href=\"" + linkUrl + "\">" + linkUrl + "</a>"

		} else {
			response = "Url add Handler - no form parameter url provided"
		}

		w.Write([]byte(response))
	}
}

func UrlRedirectIndex(c *config.AppContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		http.Redirect(w, r, "/static/index.html", http.StatusPermanentRedirect)
		return
	}
}
