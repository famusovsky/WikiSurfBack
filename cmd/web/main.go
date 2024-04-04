package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/famusovsky/WikiSurfBack/internal/app"
	"github.com/famusovsky/WikiSurfBack/internal/postgres"
	"github.com/famusovsky/WikiSurfBack/pkg/database"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
)

// TODO add names to tours

func main() {
	addr := flag.String("addr", ":8080", "HTTP address")
	overrideTables := flag.Bool("override_tables", false, "Override tables in database")
	dsn := flag.String("dsn", "", "dsn for the db")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERR\t", log.Ldate|log.Ltime)

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

	DbHandler, err := postgres.Get(db, *overrideTables)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := app.CreateApp(DbHandler, infoLog, errorLog)

	sigQuit := make(chan os.Signal, 2)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)
	eg := new(errgroup.Group)

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			return fmt.Errorf("captured signal: %v", s)
		}
	})

	app.Run(*addr)

	if err := eg.Wait(); err != nil {
		infoLog.Printf("gracefully shutting down the server: %v\n", err)
		app.Shutdown()
	}
}
