package router

import (
    "net/http"
)

func NewRouter(mux *http.ServeMux) *http.ServeMux {
    mux.HandleFunc("/", home)
    mux.HandleFunc("snippet", showSnippet)
    mux.HandleFunc("snippet/create", createSnippet)

    return mux
}


func home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    w.Write([]byte("Hello from Snippetbox"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display a specific snippet..."))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        w.Header().Set("Allow", "POST")
        w.WriteHeader(405)
        w.Write([]byte("Method not allowed"))
        return
    }

    w.Write([]byte("Create a snippet..."))
}

