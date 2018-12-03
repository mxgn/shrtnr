package handlers

import (
	"net/http"

	"github.com/mxgn/seelog"
	"github.com/mxgn/url-shrtnr/app/config"
)

var err error

var log seelog.LoggerInterface

func Init(c *config.AppContext) {
	log = c.Log
}

func UrlRedirect(c *config.AppContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Debug("GOT NEW REQUEST: Req.Host=%s, Req.URL=%s\n", r.Host, r.URL)

		code := r.URL.Path[len("/"):]

		url, err := c.DB.GetLongUrl(code)
		if err != nil {
			log.Info("UrlRedirect, short code not found:", code)
			http.Redirect(w, r, "URL "+code+" NotFound", http.StatusNotFound)
			return
		}

		log.Info("Long url from database:", url)
		http.Redirect(w, r, "http://"+string(url), http.StatusSeeOther)
	}
}

func UrlAdd(c *config.AppContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Info(`r.PostFormValue("url"):`, r.PostFormValue("url"))

		var response string

		if url := r.PostFormValue("url"); url != "" {

			shortUrl := ""
			shortUrl, err = c.DB.AddLongUrl(url)
			if err != nil {
				log.Critical("CANNOT ADD URL:", err)
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
