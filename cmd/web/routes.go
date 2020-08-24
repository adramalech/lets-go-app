package main

import "net/http"

func (app *application) routes(staticDir string) *http.ServeMux {
    fileServer := http.FileServer(http.Dir(staticDir))

    mux := http.NewServeMux()

    mux.HandleFunc("/", app.home)
    mux.HandleFunc("/snippet", app.showSnippet)
    mux.HandleFunc("/snippet/create", app.createSnippet)
    
    mux.Handle("/static/", http.StripPrefix("/static", fileServer))
    
    return mux
}
