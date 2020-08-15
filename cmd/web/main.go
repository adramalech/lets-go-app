package main

import (
    "flag"
    "log"
    "net/http"
    "os"
)

type Config struct {
    Addr string
    StaticDir string
}

func main() {
    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.LUTC)
    errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.LUTC|log.Llongfile)

    cfg := new(Config)

    flag.StringVar(&cfg.Addr, "addr", ":4000", "Http network address")
    flag.StringVar(&cfg.StaticDir, "static-dir", "./ui-static", "Path to static assets")

    flag.Parse()

    fileServer := http.FileServer(http.Dir(cfg.StaticDir))

    mux := http.NewServeMux()

    mux.HandleFunc("/", home)
    mux.HandleFunc("/snippet", showSnippet)
    mux.HandleFunc("/snippet/create", createSnippet)
    
    mux.Handle("/static/", http.StripPrefix("/static", fileServer))

    pid := os.Getpid()

    srv := &http.Server{
        Addr: cfg.Addr,
        ErrorLog: errorLog,
        Handler: mux,
    }

    infoLog.Printf("Starting server on %s with pid %d\n", cfg.Addr, pid)

    err := srv.ListenAndServe()

    errorLog.Fatal(err)
}
