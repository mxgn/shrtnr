package application
import (
	"log"
	"os"
	"path/filepath"

	"github.com/mxgn/url-shrtnr/app/storage"
)

type AppCtx struct {
	Debug     bool
	StaticDir string
	Port      string
	DB        storage.UrlDbIface
}

func (app *AppCtx) ReadConfig() {

	app.StaticDir = os.Getenv("APP_STATIC_DIR")
	if app.StaticDir == "" {
		app.StaticDir = getPath(app) + "\\www"
	}
	if app.Debug {
		log.Println(`APP_STATIC_DIR:`, app.StaticDir)
	}

	app.Port = os.Getenv("APP_PORT")
	if app.Port == "" {
		app.Port = "80"
	}
	if app.Debug {
		log.Println(`APP_PORT:`, app.Port)
	}
}

func getPath(app *AppCtx) string {

	dir, err := filepath.Abs(".") // check how it works? how get all runtime vars?

	if err != nil {
		log.Fatal(`APP_EXEC_DIR FAILED:`, err)
	}
	if app.Debug {
		log.Println(`APP_EXEC_DIR:`, dir)
	}

	return dir
}

