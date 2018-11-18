package handlers

import (
	"log"
	"net/http"

	"github.com/mxgn/url-shrtnr/app/config"
)

func UrlRedirect(c *config.AppCtx) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if c.Debug {
			log.Printf("GOT NEW REQUEST: Req.Host=%s, Req.URL=%s\n", r.Host, r.URL)
		}

		code := r.URL.Path[len("/"):]
		if code == "" {
			http.Redirect(w, r, "/static/index.html", 303)
			return
		}

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

func UrlAdd(app *config.AppCtx) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		test := r.URL.Query().Get("url")
		log.Println("r.URL.QueryGet(\"url\"):", test)

		var response string

		if url := r.PostFormValue("url"); url != "" {

			shortResult, err := app.DB.GetLongUrl(url)
			if err != nil {
				log.Panicln(`URL ADD ERROR:`, err)
			}
			linkUrl := "http://localhost/" + shortResult
			response = "<a href=\"" + linkUrl + "\">" + linkUrl + "</a>"

		} else {
			response = "Url add Handler - no form parameter url provided"
		}

		w.Write([]byte(response))
	}
}
