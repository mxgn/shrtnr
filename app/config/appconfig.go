package config

import (
	"os"
	"path/filepath"

	"github.com/mxgn/seelog"
	"github.com/mxgn/url-shrtnr/app/storage"
)

type AppContext struct {
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

func (c *AppContext) Init() {

	log = c.Log

	c.StaticDir = os.Getenv("APP_STATIC_DIR")
	if c.StaticDir == "" {
		c.StaticDir = getPath(c) + "\\www"
	}
	log.Debug(`APP_STATIC_DIR:`, c.StaticDir)

	c.Port = os.Getenv("APP_PORT")
	if c.Port == "" {
		c.Port = "80"
	}
	log.Debug(`APP_PORT:`, c.Port)

	//Postgre DB Settings init:
	c.DBcfg.Host = os.Getenv("APP_PG_HOST")
	if c.DBcfg.Host == "" {
		c.DBcfg.Host = "localhost"
	}
	log.Debug(`APP_PG_HOST: `, c.DBcfg.Host)

	c.DBcfg.Port = os.Getenv("APP_PG_PORT")
	if c.DBcfg.Port == "" {
		c.DBcfg.Port = "5432"
	}
	log.Debug(`APP_PG_PORT: `, c.DBcfg.Port)

	c.DBcfg.User = os.Getenv("APP_PG_USER")
	if c.DBcfg.User == "" {
		c.DBcfg.User = "postgres"
	}
	log.Debug(`APP_PG_USER: `, c.DBcfg.User)

	c.DBcfg.Pass = os.Getenv("APP_PG_PASS")
	log.Debug(`APP_PG_PASS: `, c.DBcfg.Pass)

	c.DBcfg.Name = os.Getenv("APP_PG_DBNAME")
	if c.DBcfg.Name == "" {
		c.DBcfg.Name = c.DBcfg.User
	}
	log.Debug(`APP_PG_DBNAME: `, c.DBcfg.Name)

}

func getPath(с *AppContext) string {

	dir, err := filepath.Abs(".") // check how it works? how get all runtime vars?
	if err != nil {
		log.Critical(`APP_EXEC_DIR FAILED:`, err)
	}

	log.Debug(`APP_EXEC_DIR:`, dir)

	return dir
}
