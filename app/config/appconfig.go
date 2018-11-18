package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mxgn/url-shrtnr/app/storage"
)

type AppCtx struct {
	Debug     bool
	StaticDir string
	Port      string
	DB        storage.UrlDbIface
	DBcfg     DBcfg
}

type DBcfg struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

func (app *AppCtx) ReadConfig() {

	app.StaticDir = os.Getenv("APP_STATIC_DIR")
	if app.StaticDir == "" {
		app.StaticDir = getPath(app) + "\\www"
	}
	if app.Debug {
		log.Println(`APP_STATIC_DIR:`, app.StaticDir)
	}

	app.Port = os.Getenv("APP_PORT")
	if app.Port == "" {
		app.Port = "80"
	}
	if app.Debug {
		log.Println(`APP_PORT:`, app.Port)
	}

	//Postgre DB Settings init:

	app.DBcfg.Host = os.Getenv("APP_PG_HOST")
	if app.DBcfg.Host == "" {
		app.DBcfg.Host = "localhost"
	}
	if app.Debug {
		log.Println(`APP_PG_HOST: `, app.DBcfg.Host)
	}

	app.DBcfg.Port = os.Getenv("APP_PG_PORT")
	if app.DBcfg.Port == "" {
		app.DBcfg.Port = "5432"
	}
	if app.Debug {
		log.Println(`APP_PG_PORT: `, app.DBcfg.Port)
	}

	app.DBcfg.User = os.Getenv("APP_PG_USER")
	if app.DBcfg.User == "" {
		app.DBcfg.User = "postgres"
	}
	if app.Debug {
		log.Println(`APP_PG_USER: `, app.DBcfg.User)
	}

	app.DBcfg.Pass = os.Getenv("APP_PG_PASS")
	if app.Debug {
		log.Println(`APP_PG_PASS: `, app.DBcfg.Pass)
	}

	app.DBcfg.Name = os.Getenv("APP_PG_DBNAME")
	if app.DBcfg.Name == "" {
		app.DBcfg.Name = app.DBcfg.User
	}
	if app.Debug {
		log.Println(`APP_PG_DBNAME: `, app.DBcfg.Name)
	}
}

func getPath(app *AppCtx) string {

	dir, err := filepath.Abs(".") // check how it works? how get all runtime vars?

	if err != nil {
		log.Fatal(`APP_EXEC_DIR FAILED:`, err)
	}
	if app.Debug {
		log.Println(`APP_EXEC_DIR:`, dir)
	}

	return dir
}
