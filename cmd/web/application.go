package main

import (
	"log"

	"github.com/adramalech/lets-go-app/snippetbox/pkg/models/mysql"
)

type application struct {
    errorLog *log.Logger
    infoLog *log.Logger
    snippets *mysql.SnippetModel
}
