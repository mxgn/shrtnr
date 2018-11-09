package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/mxgn/url-shrtnr/app/algorithm"
)

type Pgdb struct {
	Db *sql.DB
}

type Config struct {
	Host   string
	Port   string
	User   string
	Pass   string
	Dbname string
}

func (r *Pgdb) Init(cfg Config) {

	db, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.User, cfg.Pass, cfg.Dbname, cfg.Host, cfg.Port))
	if err != nil {
		log.Fatalln(err)
	}

	r.Db = db

	if err = r.Db.Ping(); err != nil {
		log.Fatalln(err)
	}
}

func (r *Pgdb) CreateSchema() {
	r.Db.Exec(`DROP TABLE URL_TBL`)
	if _, err := r.Db.Exec(`
	    CREATE TABLE IF NOT EXISTS URL_TBL (
		id    serial UNIQUE NOT NULL,
		short text   UNIQUE NOT NULL,
		url   text   UNIQUE NOT NULL
	)`); err != nil {
		log.Fatalln("!!!: ", err)
	}
}

func (r *Pgdb) GetNextId() int64 {
	stmt := `select nextval(pg_get_serial_sequence('url_tbl', 'id')) as nextId;`
	var id int64
	if err := r.Db.QueryRow(stmt).Scan(&id); err != nil {
		log.Fatalln("ERR GETING LAST ID: ", err)
	}
	log.Printf("GOT LAST ID: %v", id)

	return id
}

func (r *Pgdb) Code() string { return "nil" }

func (r *Pgdb) Save(longUrl string) string {

	stmt := `INSERT INTO URL_TBL (id, short, url) VALUES ($1, $2, $3)`

	id := r.GetNextId()
	short := algorithm.Encode(id)

	res, err := r.Db.Exec(stmt, id, short, longUrl)
	if err != nil {
		log.Fatalln("INSERT ERROR: ", err, "\n", res, "\n", id, "\n", short, "\n", longUrl)
	}
	log.Println(res)

	return "ok"
}

func (r *Pgdb) Load(string) (string, error) { return "", nil }
