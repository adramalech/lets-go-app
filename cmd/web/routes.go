package main

import (
    "net/http"

    "github.com/gorilla/mux"
)

func (app *application) addRoutes(staticDir string) *mux.Router {
    fileServer := http.FileServer(http.Dir(staticDir))

    router := mux.NewRouter()
    
    router.HandleFunc("/", app.home)
    router.HandleFunc("/snippet/{id:[0-9]+}", app.showSnippet)
    router.HandleFunc("/snippet/create", app.createSnippet).Methods("POST")
    router.HandleFunc("/snippet/create", app.createSnippetForm).Methods("GET")
    
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))
    
    return router
}
