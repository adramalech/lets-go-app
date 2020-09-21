package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/adramalech/lets-go-app/snippetbox/pkg/logger"
	"github.com/adramalech/lets-go-app/snippetbox/pkg/models/mysql"

    "github.com/justinas/alice"
)

type Config struct {
    Addr string
    StaticDir string
}

func main() {
    cfg := new(Config)
    
    ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
    defer cancel()
    
    dsn := flag.String("dsn", "web:password12345!@127.0.0.1:3306/snippetbox?parseTime=true", "MySQL data source name")
    flag.StringVar(&cfg.Addr, "addr", ":80", "Http network address")
    flag.StringVar(&cfg.StaticDir, "static-dir", "./ui-static", "Path to static assets")

    flag.Parse()

    snippetModel, dbErr := mysql.NewSnippetModel(ctx, *dsn)
    
    zLog, err := logger.NewLogger(logger.Configuration{UseJSONFormat: false})
    
    if err != nil {
        return
    }

    defer zLog.Close()

    if dbErr != nil {
        zLog.Error("Error in connecting to DB make sure DB is online and able to accept connections!")
        zLog.Fatal(err)
        return
    }

    templateCache, err := newTemplateCache("./ui/html/")

    if err != nil {
        zLog.Fatalf("error creating the template cache: %v", err)
    }

    defer snippetModel.Close()

    app := &application{
        log: zLog,
        snippets: snippetModel,
        templateCache: templateCache,
    }

    standardMiddleware := alice.New(app.recoverPanic, secureHeaders, app.logHandler, cancelHandler)

    mux := app.routes(cfg.StaticDir)
    
    reqLoggerMux := standardMiddleware.Then(mux)

    pid := os.Getpid()

    srv := &http.Server{
        Addr: cfg.Addr,
        Handler: reqLoggerMux,
    }
    
    zLog.Infof("Starting server on %s with pid %d\n", cfg.Addr, pid)

    err = srv.ListenAndServe()

    zLog.Fatal(err)
}
