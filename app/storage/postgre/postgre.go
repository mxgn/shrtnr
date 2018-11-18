package postgre

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/mxgn/url-shrtnr/app/config"
	"github.com/mxgn/url-shrtnr/app/helpers"
)

var DB *sql.DB
var err error
var debug bool

type UrlDbIface struct{}

func Init(ctx *config.AppCtx) *UrlDbIface {

	cfg := ctx.DBcfg
	debug = ctx.Debug

	DB, err = sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.User, cfg.Pass, cfg.Name, cfg.Host, cfg.Port))
	if err != nil {
		log.Fatalln(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalln(err)
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
		log.Fatalln("URL table create error:", err)
	}
}

func getNextId() int64 {
	stmt := `
			select nextval(pg_get_serial_sequence('url_tbl', 'id')) as nextId
			`
	var id int64
	if err := DB.QueryRow(stmt).Scan(&id); debug && err != nil {
		log.Println("Error getting next Id: ", err)
	}
	log.Println("Got next id:", id)
	return id
}

func checkUrl(longUrl string) string {

	var short string
	stmt := `
			SELECT short_url FROM url_tbl WHERE long_url = $1
			`

	if err := DB.QueryRow(stmt, longUrl).Scan(&short); debug && err != nil {
		log.Println("DB.QueryRow err: ", err)
	}

	if short != "" {
		log.Println("Url \"", longUrl, "\" exists, key:", short)
		return short
	}
	return ""
}

func (s UrlDbIface) AddLongUrl(longUrl string) (string, error) {

	stmt := `
			INSERT INTO URL_TBL (id, short_url, long_url) VALUES ($1, $2, $3)
			`

	if short := checkUrl(longUrl); short != "" {
		return short, nil
	}

	id := getNextId()
	short := helpers.Encode(id)

	res, err := DB.Exec(stmt, id, short, longUrl)
	if err != nil {
		log.Println("Insert error:", err)
	}
	if debug {
		log.Println("Insert result:", res)
	}
	return short, nil
}

func (s *UrlDbIface) GetLongUrl(shortUrl string) (string, error) {

	long := ""
	stmt := `SELECT long_url FROM url_tbl WHERE short_url = $1`

	if err := DB.QueryRow(stmt, shortUrl).Scan(&long); debug && err != nil {
		fmt.Println("DB SEARCH RESULT:", err)
	}

	if long == "" {
		return "", errors.New("Short url " + shortUrl + " doesnt exists")
	}
	return long, nil
}
