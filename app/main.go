package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mxgn/url-shrtnr/app/handlers"
	"github.com/mxgn/url-shrtnr/app/storages"
	"github.com/spf13/viper"
)

func main() {
	log.SetFlags(log.Lshortfile &^ (log.Ldate | log.Ltime))

	// if os.Getenv("ENVIRONMENT") == "DEV" {
	// 	viper.SetConfigName("config")
	// 	viper.SetConfigType("toml")
	// 	viper.AddConfigPath(filepath.Dir("/config"))
	// 	viper.ReadInConfig()
	// } else {
	// 	viper.AutomaticEnv()
	// }

	fmt.Println(viper.AllSettings())

	storage := &storages.Pgdb{}
	storage.Init()

	// storage.CreateSchema()
	// storage.Save("test112")

	fs := http.FileServer(http.Dir("/var/www"))
	http.Handle("/add/", http.StripPrefix("/add/", fs))

	http.Handle("/", handlers.RedirectHandler(storage))
	http.Handle("/enc/", handlers.EncodeHandler(storage))
	http.Handle("/favicon.ico", handlers.NullHandler(storage))

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
