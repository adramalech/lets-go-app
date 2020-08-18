package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "github.com/adramalech/lets-go-app/snippetbox/cmd/config"
    "strconv"
)

func Home(app *config.Application) http.HandlerFunc {
    return func (w http.ResponseWriter, r *http.Request) {
        app.InfoLog.Println("got to home!")

        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }

        files := []string{
            "./ui/html/home.page.tmpl",
            "./ui/html/base.layout.tmpl",
            "./ui/html/footer.partial.tmpl",
        }

        ts, err := template.ParseFiles(files...)

        if err != nil {
            app.ErrorLog.Println(err.Error())
            http.Error(w, "Internal Server Error", 500)
            return
        }

        err = ts.Execute(w, nil)

        if err != nil {
            log.Println(err.Error())
            http.Error(w, "Internal Server Error", 500)
        }

        w.Write([]byte("Hello from Snippetbox"))
    }
}

func ShowSnippet(app *config.Application) http.HandlerFunc { 
    return func (w http.ResponseWriter, r *http.Request) {
        app.InfoLog.Println("got to showSnippet!")

        id, err := strconv.Atoi(r.URL.Query().Get("id"))

        log.Printf("id = %v\n", id)

        if err != nil || id < 1 {
            http.NotFound(w, r)
            return
        }

        fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
    }
}

func CreateSnippet(app *config.Application) http.HandlerFunc {
    return func (w http.ResponseWriter, r *http.Request) {
        app.InfoLog.Println("got to create snippet!")

        if r.Method != "POST" {
            w.Header().Set("Allow", "POST")
            http.Error(w, "Method not allowed", 405)
            return
        }

        w.Write([]byte("Create a snippet..."))
    }
}

