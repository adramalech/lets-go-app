package main

import (
    "database/sql"
    "flag"
    "log"
    "net/http"
    "os"

    _ "github.com/go-sql-driver/mysql"
)

type Config struct {
    Addr string
    StaticDir string
}

func main() {
    cfg := new(Config)
    
    dsn := flag.String("dsn", "root:password12345@/snippetbox?parseTime=true", "MySQL data source name")

    app := &application{
        infoLog: log.New(os.Stdout, "INFO\t", log.Ldate|log.LUTC),
        errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.LUTC|log.Llongfile),
    }

    flag.StringVar(&cfg.Addr, "addr", ":4000", "Http network address")
    flag.StringVar(&cfg.StaticDir, "static-dir", "./ui-static", "Path to static assets")

    flag.Parse()

    db, err := openDB(*dsn)
    
    if err != nil {
        app.errorLog.Fatal(err)
    }

    defer db.Close()

    mux := app.routes(cfg.StaticDir)

    pid := os.Getpid()

    srv := &http.Server{
        Addr: cfg.Addr,
        ErrorLog: app.errorLog,
        Handler: mux,
    }

    app.infoLog.Printf("Starting server on %s with pid %d\n", cfg.Addr, pid)

    err = srv.ListenAndServe()

    app.errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn)
    
    if err != nil {
        return nil, err
    }
    
    err = db.Ping()

    if err != nil {
        return nil, err
    }

    return db, nil
}
