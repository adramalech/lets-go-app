package main

import (
	"context"
	"fmt"
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

    snippets, err := app.snippets.Latest(ctx)

    if err != nil {
        app.log.Error("Unable to retrieve latest snippets from database.\n")
        app.serverError(w, err)
        return
    }
    
    s := &templateData{Snippets: snippets}
    
    app.render(w, r, "home.page.tmpl", s)
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
    
    app.render(w, r, "show.page.tmpl", s)
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
    
    snip := &models.Snip{
        Content: "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\nKobayashi Issa",
        Expires: 7,
        Title: "O snail",
    }

    id, err := app.snippets.Insert(ctx, snip)

    if err != nil {
        app.log.Error("An error occurred in inserting snippet into database.\n")
        app.serverError(w, err)
        return
    }
 
    http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
