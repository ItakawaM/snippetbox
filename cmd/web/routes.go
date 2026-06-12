package main

import (
	"net/http"
	"path/filepath"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	file, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := file.Close()
			if closeErr != nil {
				return nil, err
			}

			return nil, err
		}
	}

	return file, nil
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileserver := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("GET /static/", http.StripPrefix("/static", fileserver))

	mux.HandleFunc("GET /", app.home)
	// mux.HandleFunc("POST /snippet/create", app.snippetCreate)
	mux.HandleFunc("GET /snippet/view", app.snippetView)

	return app.logRequest(secureHeaders(mux))
}
