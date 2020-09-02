package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/adramalech/lets-go-app/snippetbox/pkg/logger"
	"github.com/adramalech/lets-go-app/snippetbox/pkg/models/mysql"
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

    snippetModel, dbErr := mysql.NewSnippetModel(ctx, *dsn)
    
    zLog, err := logger.NewLogger(logger.Configuration{UseJSONFormat: false})
    
    if err != nil {
        return
    }

    defer zLog.Close()

    if dbErr != nil {
        zLog.Fatal(err)
        return
    }

    defer snippetModel.Close()

    app := &application{
        log: zLog,
        snippets: snippetModel,
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
