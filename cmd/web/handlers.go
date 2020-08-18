package main

import (
    "fmt"
    "html/template"
    "net/http"
    "strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    app.infoLog.Println("got to home!")

    if r.URL.Path != "/" {
        app.notFound(w)
        return
    }

    files := []string{
        "./ui/html/home.page.tmpl",
        "./ui/html/base.layout.tmpl",
        "./ui/html/footer.partial.tmpl",
    }

    ts, err := template.ParseFiles(files...)

    if err != nil {
        app.serverError(w, err)
        return
    }

    err = ts.Execute(w, nil)

    if err != nil {
        app.serverError(w, err)
        return
    }

    w.Write([]byte("Hello from Snippetbox"))
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
    app.infoLog.Println("got to showSnippet!")

    id, err := strconv.Atoi(r.URL.Query().Get("id"))

    app.infoLog.Printf("id = %v\n", id)

    if err != nil || id < 1 {
        app.notFound(w)
        return
    }

    fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
    app.infoLog.Println("got to create snippet!")

    if r.Method != "POST" {
        w.Header().Set("Allow", "POST")

        app.clientError(w, http.StatusMethodNotAllowed)
        return
    }

    w.Write([]byte("Create a snippet..."))
}
