package main

import (
	"log"
	"net/http"

	"github.com/mxgn/url-shrtnr/app/storages"
)

type AppHandler struct {
	storage *storages.IStorage
	H       func(*storages.IStorage, http.ResponseWriter, *http.Request) (int, error)
}

// func (ah AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	// Updated to pass ah.appContext as a parameter to our handler type.
// 	status, err := ah.H(ah.storage, w, r)
// 	if err != nil {
// 		log.Printf("HTTP %d: %q", status, err)
// 		switch status {
// 		case http.StatusNotFound:
// 			http.NotFound(w, r)
// 			// And if we wanted a friendlier error page, we can
// 			// now leverage our context instance - e.g.
// 			// err := ah.renderTemplate(w, "http_404.tmpl", nil)
// 		case http.StatusInternalServerError:
// 			http.Error(w, http.StatusText(status), status)
// 		default:
// 			http.Error(w, http.StatusText(status), status)
// 		}
// 	}
// }

// func IndexHandler(a *AppConfig, w http.ResponseWriter, r *http.Request) (int, error) {
// 	// Our handlers now have access to the members of our context struct.
// 	// e.g. we can call methods on our DB type via err := a.db.GetPosts()

// 	// fmt.Fprintf(w, "IndexHandler: db is %q and store is %q", a.db, a.store)
// 	return 200, nil
// }

func UrlRedirect(c *AppConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if c.Debug {
			log.Printf("GOT NEW REQUEST: Req.Host=%s, Req.URL=%s\n", r.Host, r.URL)
		}

		code := r.URL.Path[len("/"):]
		if code == "" {
			http.Redirect(w, r, "/static/index.html", 303)
			return
		}

		url, err := c.Storage.GetLongUrl(code)
		if err != nil {
			log.Println("")
		}

		log.Println("Long url from database:", url)
		http.Redirect(w, r, "http://"+string(url), http.StatusSeeOther)
	}
}

func UrlAdd(c *AppConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		test := r.URL.Query().Get("url")
		log.Println("r.URL.QueryGet(\"url\"):", test)

		var response string

		if url := r.PostFormValue("url"); url != "" {

			shortResult, _ := c.Storage.GetLongUrl(url)
			linkUrl := "http://localhost/" + shortResult
			response = "<a href=\"" + linkUrl + "\">" + linkUrl + "</a>"

		} else {
			response = "Url add Handler - no form parameter url provided"
		}

		w.Write([]byte(response))
	}
}
