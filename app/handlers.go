package main

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/mocdaniel/guestbook-app/app/internal/data"
)

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) createEntryHandler(w http.ResponseWriter, r *http.Request) {
	hasEntry := app.sessionManager.GetBool(r.Context(), "hasEntry")

	if hasEntry {
		err := app.writeJSON(w, http.StatusCreated, envelope{"status": "failure"}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Rating      int8   `json:"rating"`
		Testimonial string `json:"testimonial"`
		LastName    string `json:"lastname"`
		FirstName   string `json:"firstname"`
		Occupation  string `json:"occupation"`
		Github      string `json:"github"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// do some validation
	if input.Rating > 5 || input.Rating < 1 || input.Testimonial == "" || input.FirstName == "" {
		err := app.writeJSON(w, http.StatusBadGateway, envelope{"status": "failure"}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	entry := &data.Entry{
		Rating:      input.Rating,
		Testimonial: input.Testimonial,
		LastName:    input.LastName,
		FirstName:   input.FirstName,
		Occupation:  input.Occupation,
		Github:      input.Github,
	}

	err = app.models.Entries.Insert(entry)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", "/")

	app.sessionManager.Put(r.Context(), "hasEntry", true)

	err = app.writeJSON(w, http.StatusCreated, envelope{"status": "success"}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) entryHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := app.models.Entries.GetAll()
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"entries": entries}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) serveVueJS(w http.ResponseWriter, r *http.Request) {
	content, err := fs.ReadFile(app.frontend, "frontend/dist/index.html")
	if err != nil {
		app.logger.Fatal(err)
	}
	w.Write(content)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) notAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource.", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found."
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "The server encountered a problem and could not process your request."
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}
