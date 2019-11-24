package main

import (
	"database/sql"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.elastic.co/apm/module/apmgin"
	"go.elastic.co/apm/module/apmhttp"
	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	listenAddr = flag.String("listen", ":8000", "Address on which to listen for HTTP requests")
	database   = flag.String("db", "mysql:root:password@tcp(localhost:3306)/db_example?charset=utf8", "Database URL")
)

func main() {
	flag.Parse()
	// Instrument the default HTTP transport, so that outgoing
	// (reverse-proxy) requests are reported as spans.
	http.DefaultTransport = apmhttp.WrapRoundTripper(http.DefaultTransport)

	if err := Main(); err != nil {
		log.Fatal(err)
	}
}

func Main() error {

	db, err := newDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	err = initDatabase(db)
	if err != nil {
		return err
	}

	errlogfile, _ := os.Create("error.log")
	gin.DefaultErrorWriter = errlogfile

	r := gin.New()
	// Instrument gin
	r.Use(apmgin.Middleware(r))
	apiGroup := r.Group("/")
	addAPIHandlers(apiGroup, db)

	return r.Run(*listenAddr)
}

func newDatabase() (*sql.DB, error) {
	fields := strings.SplitN(*database, ":", 2)
	if len(fields) != 2 {
		return nil, errors.Errorf(
			"expected database URL with format %q, got %q",
			"<driver>:<connection-string>",
			*database,
		)
	}
	driver := fields[0]
	db, err := apmsql.Open(driver, fields[1])
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func initDatabase(db *sql.DB) error {
	file, err := ioutil.ReadFile("./db/sql/schema_mysql.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(file))
	if err != nil {
		return err
	}
	return nil

}
