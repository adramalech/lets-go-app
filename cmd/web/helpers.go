package main

import (
    "bytes"
    "fmt"
    "net/http"
    "runtime/debug"
    "time"
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

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
    if td == nil {
        td = &templateData{}
    }

    td.CurrentYear = time.Now().Year()

    return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
    templateSet, ok := app.templateCache[name]

    if !ok {
        app.serverError(w, fmt.Errorf("The template %s does not exist", name))
        return
    }

    buf := new(bytes.Buffer)

    err := templateSet.Execute(buf, td)

    if err != nil {
        app.serverError(w, err)
        return
    }

    buf.WriteTo(w)
}


