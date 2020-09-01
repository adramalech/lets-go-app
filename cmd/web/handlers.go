package main

import (
	"context"
	"fmt"
	// "html/template"
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
        app.notFound(w)
        return
    }

    // files := []string{
    //    "./ui/html/home.page.tmpl",
    //    "./ui/html/base.layout.tmpl",
    //    "./ui/html/footer.partial.tmpl",
    // }

    // ts, err := template.ParseFiles(files...)

    snippets, err := app.snippets.Latest(ctx)

    if err != nil {
        app.serverError(w, err)
        return
    }

    // err = ts.Execute(w, nil)

    // if err != nil {
    //     app.serverError(w, err)
    //    return
    // }
    
    for _, snippet := range snippets {
        fmt.Fprintf(w, "%v\n\n", snippet)
    }

    // w.Write([]byte("Hello from Snippetbox"))
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), time.Duration(60 * time.Second))
    defer cancel()
    r = r.WithContext(ctx)
    
    id, err := strconv.Atoi(r.URL.Query().Get("id"))

    app.infoLog.Printf("id = %v\n", id)

    if err != nil || id < 1 {
        app.notFound(w)
        return
    }
    
    snippet, err := app.snippets.Get(ctx, id)
    
    if err == models.ErrNoRecord {
        app.notFound(w)
        return
    } else if err != nil {
        app.serverError(w, err)
        return
    }

    fmt.Fprintf(w, "%v", snippet)
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
        app.serverError(w, err)
        return
    }
 
    http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
