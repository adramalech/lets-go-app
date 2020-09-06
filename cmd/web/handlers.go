package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
    "time"

	"github.com/adramalech/lets-go-app/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), time.Duration(60 * time.Second))
    defer cancel()
    r = r.WithContext(ctx)

    if r.URL.Path != "/" {
        app.log.Errorf("Incorrect uri provided unable to find route that matches, %s\n", r.URL.String())
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
        app.log.Error("Unable to parse html templates")
        app.serverError(w, err)
        return
    }

    snippets, err := app.snippets.Latest(ctx)

    if err != nil {
        app.log.Error("Unable to retrieve latest snippets from database.\n")
        app.serverError(w, err)
        return
    }

    s := &templateData{Snippets: snippets}

    err = ts.Execute(w, s)

    if err != nil {
        app.log.Error("Unable to render html templates")
        app.serverError(w, err)
        return
    }
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), time.Duration(60 * time.Second))
    defer cancel()
    r = r.WithContext(ctx)
    
    id, err := strconv.Atoi(r.URL.Query().Get("id"))

    if err != nil || id < 1 {
        app.log.Errorf("Id is not a correct value %d\n", id)
        app.notFound(w)
        return
    }
    
    snippet, err := app.snippets.Get(ctx, id)
    
    if err == models.ErrNoRecord {
        app.log.Errorf("No records found from id %d\n", id)
        app.notFound(w)
        return
    } else if err != nil {
        app.log.Errorf("An error occurred in getting the snippet id %d\n", id)
        app.serverError(w, err)
        return
    }

    s := &templateData{Snippet: snippet}

    files := []string{
        "./ui/html/show.page.tmpl",
        "./ui/html/base.layout.tmpl",
        "./ui/html/footer.partial.tmpl",
    }

    ts, err := template.ParseFiles(files...)
    
    if err != nil {
        app.log.Error("Unable to parse html template files.")
        app.serverError(w, err)
        return
    }
    
    err = ts.Execute(w, s)

    if err != nil {
        app.serverError(w, err)
    }
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), time.Duration(60 * time.Second))
    defer cancel()
    r = r.WithContext(ctx)

    if r.Method != "POST" {
        w.Header().Set("Allow", "POST")
        app.clientError(w, http.StatusMethodNotAllowed)
        return
    }
    
    snip := &models.Snip{}
    snip.Content = "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\nKobayashi Issa"
    snip.Expires = 7
    snip.Title = "O snail"

    id, err := app.snippets.Insert(ctx, snip)

    if err != nil {
        app.log.Error("An error occurred in inserting snippet into database.\n")
        app.serverError(w, err)
        return
    }
 
    http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
