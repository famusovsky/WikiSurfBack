package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/famusovsky/WikiSurfBack/internal/app"
	"github.com/famusovsky/WikiSurfBack/internal/postgres"
	"github.com/famusovsky/WikiSurfBack/pkg/database"
	_ "github.com/lib/pq"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP address")
	createTables := flag.Bool("create_tables", false, "Create tables in database")
	dsn := flag.String("dsn", "", "dsn for the db")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERR\t", log.Ldate|log.Ltime)

	*dsn = "postgres://postgres:qwerty@localhost:8888/postgres?sslmode=disable"

	var db *sql.DB
	var err error
	if *dsn == "" {
		db, err = database.OpenViaEnvVars("postgres")
	} else {
		db, err = database.OpenViaDsn(*dsn, "postgres")
	}
	if err != nil {
		errorLog.Fatal("error while connecting to the db ", err)
	}
	defer db.Close()

	DbHandler, err := postgres.Get(db, *createTables)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := app.CreateApp(DbHandler, infoLog, errorLog)

	app.Run(*addr)
}
