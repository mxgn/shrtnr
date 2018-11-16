package storages

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/mxgn/url-shrtnr/app/algorithm"
	"github.com/mxgn/url-shrtnr/app/models"
)

var Pgdb *DbIface

type DbIface struct {
	appCtx *models.AppConfig
	Db     *sql.DB
}

func Init(app *models.AppConfig) *DbIface {

	host := os.Getenv("APP_PG_HOST")
	if host == "" {
		host = "localhost"
	}
	if app.Debug {
		log.Println(`APP_PG_HOST: `, host)
	}

	port := os.Getenv("APP_PG_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("APP_PG_USER")
	if user == "" {
		user = "postgres"
	}
	pass := os.Getenv("APP_PG_PASS")
	if pass == "" {
		pass = ""
	}
	dbname := os.Getenv("APP_PG_DBNAME")
	if dbname == "" {
		dbname = ""
	}

	db, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, pass, dbname, host, port))
	if err != nil {
		log.Fatalln(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalln(err)
	}
	return &DbIface{Db: db}
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
