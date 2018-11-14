package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

type AppConfig struct {
	M         map[string]string
	staticDir string
	debug     bool
}

type Config interface {
	Get(key string) (string, error)
	Set(key, val string) error
}

var AppConfig AppConfig

func (c AppConfig) Get(key string) (string, error) {
	val, ok := c.M[key]
	if ok {
		return val, nil
	} else {
		return "", errors.New("Tried to get a key which doesn't exist")
	}
}

func (c AppConfig) Set(key, val string) {
	c.M[key] = val

	// debug, _ := c.debug
	if c.debug {
		log.Printf(`Set: key=%v, value=%v`, key, val)
	}

}

func Init(debug bool) (cfg AppConfig) {

	if debug {
		cfg.debug = true
		log.Println(`Init(): Staft app with debug: cfg.debug =`, debug)
	}

	cfg.M = make(map[string]string)

	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = getSelfPath(cfg) + "\\www"
	}
	cfg.Set("staticDir", staticDir)
	cfg.staticDir = staticDir

	return cfg

}

func getSelfPath(cfg AppConfig) string {
	dir, err := filepath.Abs(os.Getenv("GO_PROJECT_DIR"))
	if cfg.debug {
		log.Println(`getSelfPath(): os.Getenv("GO_PROJECT_DIR"):`, dir)
	}
	if err != nil {
		log.Fatal(`getSelfpath(): FAILED:`, err)
	}

	return dir
}
