package main

import (
    "html/template"

	"github.com/adramalech/lets-go-app/snippetbox/pkg/models/mysql"
    "github.com/adramalech/lets-go-app/snippetbox/pkg/logger"

    "github.com/gorilla/sessions"
)

type application struct {
    log logger.Logger
    snippets mysql.Snippet
    templateCache map[string]*template.Template
    session *sessions.CookieStore
}
