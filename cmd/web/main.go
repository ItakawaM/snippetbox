package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func newApplication() *application {
	return &application{
		infoLogger:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLogger: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile),
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn", "", "PostreSQL Data Source Name")
	flag.Parse()

	app := newApplication()
	db, err := openDB(*dsn)
	if err != nil {
		app.errorLogger.Fatal(err)
	}
	defer db.Close()

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: app.errorLogger,
	}

	app.infoLogger.Printf("Starting server on %s", *addr)
	app.errorLogger.Fatal(srv.ListenAndServe())
}
