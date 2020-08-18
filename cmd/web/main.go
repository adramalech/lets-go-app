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
    cfg := new(Config)
    
    app := &application{
        infoLog: log.New(os.Stdout, "INFO\t", log.Ldate|log.LUTC),
        errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.LUTC|log.Llongfile),
    }

    flag.StringVar(&cfg.Addr, "addr", ":4000", "Http network address")
    flag.StringVar(&cfg.StaticDir, "static-dir", "./ui-static", "Path to static assets")

    flag.Parse()

    mux := app.routes(cfg.StaticDir)

    pid := os.Getpid()

    srv := &http.Server{
        Addr: cfg.Addr,
        ErrorLog: app.errorLog,
        Handler: mux,
    }

    app.infoLog.Printf("Starting server on %s with pid %d\n", cfg.Addr, pid)

    err := srv.ListenAndServe()

    app.errorLog.Fatal(err)
}
