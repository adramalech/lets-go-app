package main

import "github.com/adramalech/lets-go-app/snippetbox/pkg/models"

type templateData struct {
    Snippet *models.Snippet
    Snippets []*models.Snippet
}
