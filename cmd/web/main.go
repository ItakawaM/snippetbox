package main

import (
	"log"
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
