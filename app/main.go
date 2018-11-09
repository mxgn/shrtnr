package main

import (
	"log"

	"github.com/mxgn/url-shrtnr/app/storages/postgres"
)

func main() {

	log.SetFlags(log.Lshortfile &^ (log.Ldate | log.Ltime))

	storage := &postgres.Pgdb{}
	storage.Init(postgres.Config{
		Host:   "pgdb",
		Port:   "5432",
		User:   "postgres",
		Pass:   "",
		Dbname: "postgres"})

	// storage.Db.Exec(`DROP TABLE URL_TBL`)
	// storage.CreateSchema()
	// storage.Save("test112")

	// http.Handle("/", handlers.RedirectHandler(env))
	// http.Handle("/enc/", handlers.EncodeHandler(env))
	// http.Handle("/dec/", handlers.DecodeHandler(env))

	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8080"
	// }
	// if err := http.ListenAndServe(":"+port, nil); err != nil {
	// 	log.Fatal(err)
	// }

}
