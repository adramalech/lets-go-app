package main

import (
    "github.com/adramalech/lets-go-app/snippetbox/web"
    "log"
    "net/http"
    "os"
)

func main() {
    mux := http.NewServeMux()

    mux = router.NewRouter(mux)
    
    pid := os.Getpid()

    log.Printf("Starting server on :4000 with pid %d\n", pid)

    err := http.ListenAndServe(":4000", mux)

    log.Fatal(err)
}
