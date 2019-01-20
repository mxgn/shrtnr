package storage

import (
	"errors"
	"strings"
)

type UrlDbIface interface {
	AddLongUrl(string) (string, error)
	GetLongUrl(string) (string, error)
	// Search(r UrlRec) (chan *UrlSearch, error)
}

//
//----------------
//
// NotFoundError indicates that a record could not be located.
// This differentiates between not finding a record and the
// storage layer having an error.
type NotFoundError struct {
	error
}

func (n NotFoundError) isNotFound() {}

// NotFound indicates if the error is that the ID could
// not be found.
func NotFound(e error) bool {
	if _, ok := e.(NotFoundError); ok {
		return true
	}
	return false
}

// UrlRec represents an url record.
type UrlRec struct {
	ID       uint64
	ShortUrl string
	LongUrl  string
}

// Validate validates the fields are valid.
func (e *UrlRec) Validate() error {
	if e.ID == 0 {
		return errors.New("ID field cannot be 0")
	}

	switch "" {
	case strings.TrimSpace(e.LongUrl):
		return errors.New("First field cannot be empty string")
	case strings.TrimSpace(e.ShortUrl):
		return errors.New("Last field cannot be empty string")
	}
	return nil
}

type UrlSearch struct {
	// Rec exists if a valid response was returned.
	Rec *UrlRec
	// Err exists if the storage system had an error mid search.
	Err error
}
