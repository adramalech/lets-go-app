package mysql

import (
    "context"
    "database/sql"
	
    _ "github.com/go-sql-driver/mysql"

    "github.com/jmoiron/sqlx"

    "github.com/adramalech/lets-go-app/snippetbox/pkg/models"
)

type Snippet interface {
    Insert(ctx context.Context, s *models.Snip) (int, error)
    Get(ctx context.Context, id int) (*models.Snippet, error)
    Latest(ctxt context.Context) ([]*models.Snippet, error)
    Close() error
}

type SnippetModel struct {
    DB *sqlx.DB
}

func NewSnippetModel(ctx context.Context, dsn string) (Snippet, error) {
    db := sqlx.MustOpen("mysql", dsn)
    
    err := db.PingContext(ctx)

    if err != nil {
        return nil, err
    }
    
    return &SnippetModel{DB: db}, nil
}

func (m *SnippetModel) Close() error {
    err := m.Close()
    return err
}

func (m *SnippetModel) Insert(ctx context.Context, s *models.Snip) (int, error) {
    stmt := `
        INSERT INTO snippets (title, content, created, expires)
        VALUES (:title, :content, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL :expires DAY))
    `

    result, err := m.DB.NamedExecContext(ctx, stmt, s)

    if err != nil {
        return 0, err
    }

    id, err := result.LastInsertId()

    if err != nil {
        return 0, err
    }

    return int(id), nil
}

func (m *SnippetModel) Get(ctx context.Context, id int) (*models.Snippet, error) {
    stmt := `
        SELECT id, title, content, created, expires
        FROM snippets
        WHERE expires > UTC_TIMESTAMP() AND id = ?
    `

    row := m.DB.QueryRowxContext(ctx, stmt, id)

    snippet := &models.Snippet{}
    
    err := row.StructScan(snippet)

    if err == sql.ErrNoRows {
        return nil, models.ErrNoRecord
    } else if err != nil {
        return nil, err
    }

    return snippet, nil
}

func (m *SnippetModel) Latest(ctx context.Context) ([]*models.Snippet, error) {
    stmt := `
        SELECT id, title, content, created, expires
        FROM snippets
        WHERE expires > UTC_TIMESTAMP()
        ORDER BY created DESC
        LIMIT 10
    `
    
    rows, err := m.DB.QueryxContext(ctx, stmt)

    if err == sql.ErrNoRows {
        return nil, models.ErrNoRecord
    } else if err != nil {
        return nil, err
    }

    defer rows.Close()

    snippets := []*models.Snippet{}

    for rows.Next() {
        snippet := &models.Snippet{}
        
        err := rows.StructScan(&snippet)

        if err != nil {
            return nil, err
        }

        snippets = append(snippets, snippet)
    }

    err = rows.Err()

    if err != nil {
        return nil, err
    }

    return snippets, nil
}
