package main

import (
    "fmt"
    "net/http"
    "runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
    trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
    app.log.Errorf("%v\n", trace)
    statusText := http.StatusText(http.StatusInternalServerError)
    app.log.Error("Sending response of %d with %s\n", http.StatusInternalServerError, statusText)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
    statusText := http.StatusText(status)
    app.log.Errorf("An Error occurred sending back status code %d %s\n", status, statusText)
    http.Error(w, statusText, status)
}

func (app *application) notFound(w http.ResponseWriter) {
    app.clientError(w, http.StatusNotFound)
}
