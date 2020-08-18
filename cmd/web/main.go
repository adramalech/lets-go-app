package main

import (
    "flag"
    "github.com/adramalech/lets-go-app/snippetbox/cmd/config"
    "log"
    "net/http"
    "os"
)

type Config struct {
    Addr string
    StaticDir string
}

func main() {
    cfg := new(Config)
    
    app := &config.Application{
        ErrorLog: log.New(os.Stdout, "INFO\t", log.Ldate|log.LUTC),
        InfoLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.LUTC|log.Llongfile),
    }

    flag.StringVar(&cfg.Addr, "addr", ":4000", "Http network address")
    flag.StringVar(&cfg.StaticDir, "static-dir", "./ui-static", "Path to static assets")

    flag.Parse()

    fileServer := http.FileServer(http.Dir(cfg.StaticDir))

    mux := http.NewServeMux()

    mux.HandleFunc("/", Home(app))
    mux.HandleFunc("/snippet", ShowSnippet(app))
    mux.HandleFunc("/snippet/create", CreateSnippet(app))
    
    mux.Handle("/static/", http.StripPrefix("/static", fileServer))

    pid := os.Getpid()

    srv := &http.Server{
        Addr: cfg.Addr,
        ErrorLog: app.ErrorLog,
        Handler: mux,
    }

    app.InfoLog.Printf("Starting server on %s with pid %d\n", cfg.Addr, pid)

    err := srv.ListenAndServe()

    app.ErrorLog.Fatal(err)
}
