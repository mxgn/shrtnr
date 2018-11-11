package storages

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/mxgn/url-shrtnr/app/algorithm"
)

var Pgdb DbIface

//db - local package var
var db *sql.DB

type DbIface struct {
	Db *sql.DB
}

type Config struct {
	Host   string
	Port   string
	User   string
	Pass   string
	Dbname string
}

func (r *DbIface) Init() {

	cfg := &Config{}

	cfg.Host = os.Getenv("PG_HOST")
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}
	cfg.Port = os.Getenv("PG_PORT")
	if cfg.Port == "" {
		cfg.Port = "5432"
	}
	cfg.User = os.Getenv("PG_USER")
	if cfg.User == "" {
		cfg.User = "postgres"
	}
	cfg.Pass = os.Getenv("PG_PASS")
	if cfg.Pass == "" {
		cfg.Pass = ""
	}
	cfg.Dbname = os.Getenv("PG_DBNAME")
	if cfg.Dbname == "" {
		cfg.Dbname = ""
	}

	db, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.User, cfg.Pass, cfg.Dbname, cfg.Host, cfg.Port))
	if err != nil {
		log.Fatalln(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalln(err)
	}
	Pgdb.Db = db
}

func (r *DbIface) CreateSchema() {
	// r.Db.Exec(`DROP TABLE URL_TBL`)
	stmt := `
			CREATE TABLE IF NOT EXISTS URL_TBL (
				id         serial UNIQUE NOT NULL,
				short_url  text   UNIQUE NOT NULL,
				long_url   text   UNIQUE NOT NULL
			)`
	if _, err := Pgdb.Db.Exec(stmt); err != nil {
		log.Fatalln("URL table create error:", err)
	}
}

func (r *DbIface) GetNextId() int64 {
	stmt := `
			select nextval(pg_get_serial_sequence('url_tbl', 'id')) as nextId
			`
	var id int64
	if err := r.Db.QueryRow(stmt).Scan(&id); err != nil {
		log.Println("Error getting next Id: ", err)
	}
	log.Println("Got next id:", id)
	return id
}

func (r *DbIface) checkUrl(longUrl string) string {

	var short string
	stmt := `
			SELECT short_url FROM url_tbl WHERE long_url = $1
			`

	if err := r.Db.QueryRow(stmt, longUrl).Scan(&short); err != nil {
		log.Println(err)
	}

	if short != "" {
		log.Println("Url \"", longUrl, "\" exists, key:", short)
		return short
	}
	return ""
}

func (r *DbIface) Save(longUrl string) string {

	stmt := `
			INSERT INTO URL_TBL (id, short_url, long_url) VALUES ($1, $2, $3)
			`

	if short := r.checkUrl(longUrl); short != "" {
		return short
	}

	id := r.GetNextId()
	short := algorithm.Encode(id)

	res, err := r.Db.Exec(stmt, id, short, longUrl)
	if err != nil {
		log.Println("Insert error:", err)
	}
	log.Println("Insert result:", res)

	return short
}

func (r *DbIface) Load(shortUrl string) (string, error) {

	long := ""
	stmt := `SELECT long_url FROM url_tbl WHERE short_url = $1`

	if err := r.Db.QueryRow(stmt, shortUrl).Scan(&long); err != nil {
		fmt.Println("!!!Short:", shortUrl, "\n\nERROR:", err)
	}

	if long == "" {
		return "", errors.New("Short url " + shortUrl + " doesnt exists")
	}
	return long, nil
}
