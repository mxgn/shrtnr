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

	app.StaticDir = os.Getenv("APP_STATIC_DIR")
	if app.StaticDir == "" {
		app.StaticDir = getPath(app) + "\\www"
	}
	if app.Debug {
		log.Println(`APP_STATIC_DIR:`, app.StaticDir)
	}

	app.Port = os.Getenv("APP_PORT")
	if app.Port == "" {
		app.Port = "8080"
	}
	if app.Debug {
		log.Println(`APP_PORT:`, app.Port)
	}
}

func getPath(app *AppConfig) string {
	dir, err := filepath.Abs(os.Getenv("APP_EXEC_DIR"))
	if app.Debug {
		log.Println(`APP_EXEC_DIR:`, dir)
	}
	if err != nil {
		log.Fatal(`APP_EXEC_DIR FAILED:`, err)
	}

	return dir
}
