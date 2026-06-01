package main

import (
	"flag"
	"log"
	"net/http"
	"os"
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

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	flag.Parse()

	app := newApplication()
	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: app.errorLogger,
	}

	app.infoLogger.Printf("Starting server on %s", *addr)
	app.errorLogger.Fatal(srv.ListenAndServe())
}
