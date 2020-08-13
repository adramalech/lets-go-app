package main

import (
    "log"
    "net/http"
    "os"
)

func main() {
    fileServer := http.FileServer(http.Dir("./ui/static"))

    mux := http.NewServeMux()

    mux.HandleFunc("/", home)
    mux.HandleFunc("/snippet", showSnippet)
    mux.HandleFunc("/snippet/create", createSnippet)
    
    mux.Handle("/static/", http.StripPrefix("/static", fileServer))

    pid := os.Getpid()

    log.Printf("Starting server on :4000 with pid %d\n", pid)

    err := http.ListenAndServe(":4000", mux)

    log.Fatal(err)
}
