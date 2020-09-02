package main

import (
	"github.com/adramalech/lets-go-app/snippetbox/pkg/models/mysql"
    "github.com/adramalech/lets-go-app/snippetbox/pkg/logger"
)

type application struct {
    log logger.Logger
    snippets *mysql.SnippetModel
}
