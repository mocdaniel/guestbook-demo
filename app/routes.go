package main

import (
	"io/fs"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.notAllowedResponse)

	assets, err := fs.Sub(app.frontend, "frontend/dist/assets")
	if err != nil {
		app.logger.Fatal(err)
	}

	public, err := fs.Sub(app.frontend, "frontend/dist")
	if err != nil {
		app.logger.Fatal(err)
	}

	router.HandlerFunc(http.MethodGet, "/", app.serveVueJS)
	router.HandlerFunc(http.MethodGet, "/v1/entries", app.entryHandler)
	router.HandlerFunc(http.MethodPost, "/v1/entries", app.createEntryHandler)
	router.ServeFiles("/assets/*filepath", http.FS(assets))
	router.ServeFiles("/img/*filepath", http.FS(public))

	return router
}
