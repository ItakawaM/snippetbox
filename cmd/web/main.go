package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/ItakawaM/snippetbox/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	snippets    *models.SnippetModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn", "", "PostreSQL Data Source Name")
	flag.Parse()

	infoLogger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLogger.Fatal(err)
	}
	defer db.Close()

	app := &application{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		snippets:    &models.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: app.errorLogger,
	}

	infoLogger.Printf("Starting server on %s", *addr)
	errorLogger.Fatal(srv.ListenAndServe())
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
