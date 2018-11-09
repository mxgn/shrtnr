package main

import (
	"strconv"

	"github.com/mxgn/url-shrtnr/app/storages/postgres"
)

func main() {

	// storage := &storages.Redis{}
	// if err := storage.Init(); err != nil {
	// 	log.Fatal(err)
	// }

	storage := &postgres.Pgdb{}
	storage.Init(postgres.Config{
		Host:   "localhost",
		Port:   "5432",
		User:   "postgres",
		Pass:   "",
		Dbname: "postgres"})

	storage.CreateSchema()
	for i := 0; i < 1000; i++ {
		storage.Save("str" + strconv.Itoa(i))
	}
	// storage.Save("sdfsdfsd")

	// db.Exec(`DROP TABLE URL_TBL`)

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
