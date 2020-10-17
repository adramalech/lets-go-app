package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/adramalech/lets-go-app/snippetbox/pkg/models"
    "github.com/adramalech/lets-go-app/snippetbox/pkg/forms"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
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
    ctx := r.Context()
    
    id, err := strconv.Atoi(r.URL.Query().Get(":id"))

    if err != nil && id < 1 {
        app.log.Errorf("Id needs to exist and be a correct number, id=%d", id)
        app.log.Errorf("Error: %v", err)
        app.clientError(w, http.StatusBadRequest)
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

    s := &templateData{
        Snippet: snippet,
    }
    
    app.render(w, r, "show.page.tmpl", s)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    err := r.ParseForm()

    if err != nil {
        app.clientError(w, http.StatusBadRequest)
    }

    form := forms.New(r.PostForm)

    form.Required("title", "content", "expires")
    form.MaxLength("title", 100)
    form.PermittedValues("expires", "365", "7", "1")

    if !form.Valid() {
        app.render(w, r, "create.page.tmpl", &templateData{
            Form: form,
        })

        return
    }
    
    session, err := app.session.Get(r, "session")

    if err != nil {
        app.log.Errorf("An error occurred in getting the session %v", err)
        app.serverError(w, err)
        return
    }

    expiresStr := form.Get("expires")
    
    expires, err := strconv.Atoi(expiresStr)

    if err != nil {
        app.log.Errorf("Unable to parse expires field, %s", expiresStr)
        app.serverError(w, err)
        return
    }

    snip := &models.Snip{
        Title: form.Get("title"), 
        Content: form.Get("content"), 
        Expires: expires,
    }

    id, err := app.snippets.Insert(ctx, snip)
    
    if err != nil {
        app.log.Error("An error occurred in inserting snippet into database.\n")
        app.serverError(w, err)
        return
    }

    session.AddFlash("Snippet successfully created!", "flash")
    err = session.Save(r, w)
    
    if err != nil {
        app.log.Errorf("Unable to save session flash message %v", err)
    }

    http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
    app.render(w, r, "create.page.tmpl", &templateData{
        Form: forms.New(nil),
    })
}
