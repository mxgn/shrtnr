package postgre

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mxgn/seelog"
	"github.com/mxgn/url-shrtnr/app/config"
	"github.com/mxgn/url-shrtnr/app/helpers"
)

var DB *sql.DB
var log seelog.LoggerInterface
var err error
var debug bool

type UrlDbIface struct{}

func Init(c *config.AppContext) *UrlDbIface {

	cfg := c.DBcfg
	log = c.Log

	DB, err = sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.User, cfg.Pass, cfg.Name, cfg.Host, cfg.Port))
	if err != nil {
		log.Critical(err)
	}

	if err = DB.Ping(); err != nil {
		log.Critical(err)
	}
	return &UrlDbIface{}
}

func CreateSchema() {
	DB.Exec(`DROP TABLE URL_TBL`)
	stmt := `
			CREATE TABLE IF NOT EXISTS URL_TBL (
				id         serial UNIQUE NOT NULL,
				short_url  text   UNIQUE NOT NULL,
				long_url   text   UNIQUE NOT NULL
			)`
	if _, err := DB.Exec(stmt); err != nil {
		log.Critical("URL table create error:", err)
	}
}

func getNextId() int64 {
	stmt := `
			select nextval(pg_get_serial_sequence('url_tbl', 'id')) as nextId
			`
	var id int64
	if err := DB.QueryRow(stmt).Scan(&id); debug && err != nil {
		log.Error("Error getting next Id: ", err)
	}
	log.Info("Got next id:", id)
	return id
}

func checkUrl(longUrl string) string {

	log.Debug("Entry to check url, with arg:", longUrl)

	var short string
	stmt := `
			SELECT short_url FROM url_tbl WHERE long_url = $1
			`

	log.Trace("Checking URL before add:", longUrl)
	if err := DB.QueryRow(stmt, longUrl).Scan(&short); debug && err != nil {
		log.Error("DB.QueryRow err: ", err)
	}

	if short != "" {
		log.Info("Url \"", longUrl, "\" exists, key:", short)
		return short
	}
	return ""
}

func (s UrlDbIface) AddLongUrl(longUrl string) (string, error) {

	defer helpers.Un(helpers.Trace("postgre.AddLongUrl"))

	stmt := `
			INSERT INTO URL_TBL (id, short_url, long_url) VALUES ($1, $2, $3)
			`

	if short := checkUrl(longUrl); short != "" {
		return short, nil
	}

	id := getNextId()
	short := helpers.Encode(id)

	res, err := DB.Exec(stmt, id, short, longUrl)
	if debug && err != nil {
		log.Error("Insert error:", err)
		log.Error("Insert result:", res)
	}
	return short, nil
}

func (s *UrlDbIface) GetLongUrl(shortUrl string) (string, error) {

	long := ""
	stmt := `SELECT long_url FROM url_tbl WHERE short_url = $1`

	if err := DB.QueryRow(stmt, shortUrl).Scan(&long); err != nil {
		log.Error("stmt: %v\n result: %v", stmt, err)
	}
	log.Infof("DB SEARCH RESULT: %v", long)

	if long == "" {
		return "", errors.New("Short url " + shortUrl + " doesnt exists")
	}
	return long, nil
}
