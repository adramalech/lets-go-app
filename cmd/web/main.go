package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/adramalech/lets-go-app/snippetbox/pkg/logger"
	"github.com/adramalech/lets-go-app/snippetbox/pkg/models/mysql"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
    Addr string
    StaticDir string
}

func main() {
    cfg := new(Config)
    
    ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
    defer cancel()
    
    dsn := flag.String("dsn", "root:password12345@/snippetbox?parseTime=true", "MySQL data source name")
    flag.StringVar(&cfg.Addr, "addr", ":4000", "Http network address")
    flag.StringVar(&cfg.StaticDir, "static-dir", "./ui-static", "Path to static assets")

    flag.Parse()

    db, dbErr := openDB(ctx, *dsn)
    
    zLog, err := logger.NewLogger(logger.Configuration{UseJSONFormat: false})
    
    if err != nil {
        return
    }

    defer zLog.Close()

    if dbErr != nil {
        zLog.Fatal(err)
        return
    }

    defer db.Close()

    app := &application{
        log: zLog,
        snippets: &mysql.SnippetModel{DB: db},
    }

    mux := app.routes(cfg.StaticDir)
    
    reqLoggerMux := logHandler(mux, zLog)

    pid := os.Getpid()

    srv := &http.Server{
        Addr: cfg.Addr,
        Handler: reqLoggerMux,
    }
    
    zLog.Infof("Starting server on %s with pid %d\n", cfg.Addr, pid)

    err = srv.ListenAndServe()

    zLog.Fatal(err)
}

func openDB(ctx context.Context, dsn string) (*sqlx.DB, error) {
    db := sqlx.MustOpen("mysql", dsn)
    
    err := db.PingContext(ctx)

    if err != nil {
        return nil, err
    }

    return db, nil
}
