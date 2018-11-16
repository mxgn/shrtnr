package models

import (
	"log"
	"os"
	"path/filepath"
)

type AppConfig struct {
	StaticDir string
	Debug     bool
	Port      string
}

func (app *AppConfig) Init() {

	if app.Debug {
		log.Println(`Init() - DEBUG:`, app.Debug)
	}

	app.StaticDir = os.Getenv("STATIC_DIR")
	if app.StaticDir == "" {
		app.StaticDir = getPath(app) + "\\www"
	}
	if app.Debug {
		log.Println(`Init() - STATIC_DIR:`, app.StaticDir)
	}

	app.Port = os.Getenv("PORT")
	if app.Port == "" {
		app.Port = "80"
	}
	if app.Debug {
		log.Println(`Init() - PORT:`, app.Port)
	}
}

func getPath(app *AppConfig) string {
	dir, err := filepath.Abs(os.Getenv("GO_PROJECT_DIR"))
	if app.Debug {
		log.Println(`getPath() - DIR:`, dir)
	}
	if err != nil {
		log.Fatal(`getPath() - FAILED:`, err)
	}

	return dir
}
