package helpers

import (
	"github.com/mxgn/seelog"
	"github.com/mxgn/url-shrtnr/app/config"
)

var log seelog.LoggerInterface

func Init(c *config.AppContext) {
	log = c.Log
}

func Trace(s string) string {
	log.Trace("Entering:", s)
	return s
}

func Un(s string) {
	log.Trace("Leaving:", s)
}
