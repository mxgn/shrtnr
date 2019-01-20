package postgre

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mxgn/url-shrtnr/app/config"
	"github.com/mxgn/url-shrtnr/app/helpers"
)

type dbImpl struct{}

var (
	db    *sql.DB
	err   error
	debug bool
)

func Init(ctx *config.AppContext) *dbImpl {

	cfg := c.DBcfg
	log = c.Log

	db, err = sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.User, cfg.Pass, cfg.Name, cfg.Host, cfg.Port))
	if err != nil {
		log.Critical(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalln(err)
	}
	return &dbImpl{}
}

func CreateSchema() {
	db.Exec(`DROP TABLE URL_TBL`)
	stmt := `
			CREATE TABLE IF NOT EXISTS URL_TBL (
				id         serial UNIQUE NOT NULL,
				short_url  text   UNIQUE NOT NULL,
				long_url   text   UNIQUE NOT NULL
			)`
	if _, err := db.Exec(stmt); err != nil {
		log.Fatalln("URL table create error:", err)
	}
}

func getNextId() int64 {
	stmt := `
			select nextval(pg_get_serial_sequence('url_tbl', 'id')) as nextId
			`
	var id int64
	if err := db.QueryRow(stmt).Scan(&id); debug && err != nil {
		log.Println("Error getting next Id: ", err)
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

	log.Println("Checking URL before add:", longUrl)
	if err := db.QueryRow(stmt, longUrl).Scan(&short); debug && err != nil {
		log.Println("DB.QueryRow err: ", err)
	}

	if short != "" {
		log.Info("Url \"", longUrl, "\" exists, key:", short)
		return short
	}
	return ""
}

func (s *dbImpl) AddLongUrl(longUrl string) (string, error) {

	defer helpers.Un(helpers.Trace("postgre.AddLongUrl"))

	stmt := `
			INSERT INTO URL_TBL (id, short_url, long_url) VALUES ($1, $2, $3)
			`

	if short := checkUrl(longUrl); short != "" {
		return short, nil
	}

	id := getNextId()
	short := helpers.Encode(id)

	res, err := db.Exec(stmt, id, short, longUrl)
	if debug && err != nil {
		log.Error("Insert error:", err)
		log.Error("Insert result:", res)
	}
	return short, nil
}

func (s *dbImpl) GetLongUrl(shortUrl string) (string, error) {

	long := ""
	stmt := `SELECT long_url FROM url_tbl WHERE short_url = $1`

	if err := db.QueryRow(stmt, shortUrl).Scan(&long); debug {
		fmt.Println("DB SEARCH RESULT:", long, err)
	}
	log.Infof("DB SEARCH RESULT: %v", long)

	if long == "" {
		return "", errors.New("Short url " + shortUrl + " doesnt exists")
	}
	return long, nil
}
