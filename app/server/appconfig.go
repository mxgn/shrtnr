package server

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

type AppConfig struct {
	M         map[string]string
	StaticDir string
	Debug     bool
}

var Cfg *AppConfig

type Config interface {
	Get(key string) (string, error)
	Set(key, val string) error
}

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
	if c.Debug {
		log.Printf(`Set: key=%v, value=%v`, key, val)
	}

}

func Init(cfg *AppConfig) {

	if cfg.Debug {
		cfg.Debug = true
		log.Println(`Init(): Staft app with debug: cfg.debug =`, cfg.Debug)
	}

	cfg.M = make(map[string]string)

	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = getSelfPath(cfg) + "\\www"
	}
	cfg.Set("staticDir", staticDir)
	cfg.StaticDir = staticDir
}

func getSelfPath(cfg *AppConfig) string {
	dir, err := filepath.Abs(os.Getenv("GO_PROJECT_DIR"))
	if cfg.Debug {
		log.Println(`getSelfPath(): os.Getenv("GO_PROJECT_DIR"):`, dir)
	}
	if err != nil {
		log.Fatal(`getSelfpath(): FAILED:`, err)
	}

	return dir
}
