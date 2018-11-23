package config

import (
	"os"
	"path/filepath"

	"github.com/cihub/seelog"
	"github.com/mxgn/url-shrtnr/app/storage"
	// log "github.com/sirupsen/logrus"
)

type AppContext struct {
	Debug     bool
	StaticDir string
	Port      string
	DB        storage.UrlDbIface
	DBcfg     DBcfg
	Log       seelog.LoggerInterface
}

type DBcfg struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

var log seelog.LoggerInterface

func (app *AppContext) ReadConfig() {

	log = app.Log

	app.StaticDir = os.Getenv("APP_STATIC_DIR")
	if app.StaticDir == "" {
		app.StaticDir = getPath(app) + "\\www"
	}

	log.Debug(`APP_STATIC_DIR:`, app.StaticDir)

	app.Port = os.Getenv("APP_PORT")
	if app.Port == "" {
		app.Port = "80"
	}

	log.Debug(`APP_PORT:`, app.Port)

	//Postgre DB Settings init:

	app.DBcfg.Host = os.Getenv("APP_PG_HOST")
	if app.DBcfg.Host == "" {
		app.DBcfg.Host = "localhost"
	}
	log.Warn(`APP_PG_HOST: `, app.DBcfg.Host)

	app.DBcfg.Port = os.Getenv("APP_PG_PORT")
	if app.DBcfg.Port == "" {
		app.DBcfg.Port = "5432"
	}
	log.Debug(`APP_PG_PORT: `, app.DBcfg.Port)

	app.DBcfg.User = os.Getenv("APP_PG_USER")
	if app.DBcfg.User == "" {
		app.DBcfg.User = "postgres"
	}
	log.Debug(`APP_PG_USER: `, app.DBcfg.User)

	app.DBcfg.Pass = os.Getenv("APP_PG_PASS")
	log.Debug(`APP_PG_PASS: `, app.DBcfg.Pass)

	app.DBcfg.Name = os.Getenv("APP_PG_DBNAME")
	if app.DBcfg.Name == "" {
		app.DBcfg.Name = app.DBcfg.User
	}
	log.Debug(`APP_PG_DBNAME: `, app.DBcfg.Name)

}

func getPath(app *AppContext) string {

	dir, err := filepath.Abs(".") // check how it works? how get all runtime vars?

	if err != nil {
		log.Critical(`APP_EXEC_DIR FAILED:`, err)
	}
	if app.Debug {
		log.Debug(`APP_EXEC_DIR:`, dir)
	}

	return dir
}
