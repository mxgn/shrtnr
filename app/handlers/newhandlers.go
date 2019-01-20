package handlers

import (
	"fmt"

	"net/http"

	"github.com/mxgn/url-shrtnr/app/config"
)

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

// The Handler struct that takes a configured Env and a function matching
// our useful signature.
type Handler struct {
	Ctx *config.AppContext
	H   func(e *config.AppContext, w http.ResponseWriter, r *http.Request) error
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.Ctx, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific
			// HTTP status code.
			h.Ctx.Log.Error("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

func RedirectNew(ctx *config.AppContext, w http.ResponseWriter, r *http.Request) error {

	code := r.URL.Path[len("/"):]
	if code == "" {
		http.Redirect(w, r, "/static/index.html", http.StatusPermanentRedirect)
		return nil
	}

	url, err := ctx.DB.GetLongUrl(code)
	if err != nil {
		// We return a status error here, which conveniently wraps the error
		// returned from our DB queries. We can clearly define which errors
		// are worth raising a HTTP 500 over vs. which might just be a HTTP
		// 404, 403 or 401 (as appropriate). It's also clear where our
		// handler should stop processing by returning early.
		return StatusError{500, err}
	}

	http.Redirect(w, r, "http://"+string(url), http.StatusSeeOther)
	return nil
}

func AddNew(c *config.AppContext, w http.ResponseWriter, r *http.Request) error {

	// test := r.URL.Query().Get("url")
	c.Log.Info(`Post form "url" param value:`, r.PostFormValue("url"))

	var response string

	if url := r.PostFormValue("url"); url != "" {

		shortUrl := ""
		shortUrl, err = c.DB.AddLongUrl(url)
		if err != nil {
			c.Log.Critical("CANNOT ADD URL:", err)
		}

		linkUrl := "http://localhost/" + shortUrl
		response = "<a href=\"" + linkUrl + "\">" + linkUrl + "</a>"

	} else {
		response = "No form parameter url provided"
	}
	fmt.Fprintf(w, "%+v", response)
	return nil

}

func IndexRedirectNew(ctx *config.AppContext, w http.ResponseWriter, r *http.Request) error {

	http.Redirect(w, r, "/static/index.html", http.StatusPermanentRedirect)
	return nil
}
