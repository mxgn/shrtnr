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

type IStorageImpl struct {
	db *sql.DB
}

var DB *sql.DB
var err error

func Init(debug bool) *IStorageImpl {

	host := os.Getenv("APP_PG_HOST")
	if host == "" {
		host = "localhost"
	}
	if debug {
		log.Println(`APP_PG_HOST: `, host)
	}

	port := os.Getenv("APP_PG_PORT")
	if port == "" {
		port = "5432"
	}
	if debug {
		log.Println(`APP_PG_PORT: `, port)
	}

	user := os.Getenv("APP_PG_USER")
	if user == "" {
		user = "postgres"
	}
	if debug {
		log.Println(`APP_PG_USER: `, user)
	}

	pass := os.Getenv("APP_PG_PASS")
	if pass == "" {
		pass = ""
	}
	if debug {
		log.Println(`APP_PG_PASS: `, pass)
	}

	dbname := os.Getenv("APP_PG_DBNAME")
	if dbname == "" {
		dbname = ""
	}
	if debug {
		log.Println(`APP_PG_DBNAME: `, dbname)
	}

	DB, err = sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, pass, dbname, host, port))
	if err != nil {
		log.Fatalln(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalln(err)
	}
	return &IStorageImpl{DB}
}

func сreateSchema() {
	// r.Db.Exec(`DROP TABLE URL_TBL`)
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
	if err := DB.QueryRow(stmt).Scan(&id); err != nil {
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

	if err := DB.QueryRow(stmt, longUrl).Scan(&short); err != nil {
		log.Println(err)
	}

	if short != "" {
		log.Println("Url \"", longUrl, "\" exists, key:", short)
		return short
	}
	return ""
}

func (s IStorageImpl) AddLongUrl(longUrl string) (string, error) {

	stmt := `
			INSERT INTO URL_TBL (id, short_url, long_url) VALUES ($1, $2, $3)
			`

	if short := checkUrl(longUrl); short != "" {
		return short, nil
	}

	id := getNextId()
	short := algorithm.Encode(id)

	res, err := DB.Exec(stmt, id, short, longUrl)
	if err != nil {
		log.Println("Insert error:", err)
	}
	log.Println("Insert result:", res)

	return short, nil
}

func (s *IStorageImpl) GetLongUrl(shortUrl string) (string, error) {

	long := ""
	stmt := `SELECT long_url FROM url_tbl WHERE short_url = $1`

	if err := DB.QueryRow(stmt, shortUrl).Scan(&long); err != nil {
		fmt.Println("!!!Short:", shortUrl, "\n\nERROR:", err)
	}

	if long == "" {
		return "", errors.New("Short url " + shortUrl + " doesnt exists")
	}
	return long, nil
}
