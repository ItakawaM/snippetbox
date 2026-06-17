package models

import (
	"database/sql"
	"errors"
	"time"
)

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	statement := `INSERT INTO snippets (title, content, created, expires)
VALUES ($1, $2, NOW() AT TIME ZONE 'UTC', NOW() AT TIME ZONE 'UTC' + ($3 * INTERVAL '1 day')) RETURNING id;`

	var id int
	if err := m.DB.QueryRow(statement, title, content, expires).Scan(&id); err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	statement := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > NOW() AT TIME ZONE 'UTC' AND id = $1`

	snippet := &Snippet{}
	if err := m.DB.QueryRow(statement, id).
		Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}

		return nil, err
	}

	return snippet, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	statement := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > NOW() AT TIME ZONE 'UTC' ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*Snippet{}
	for rows.Next() {
		snippet := &Snippet{}
		if err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires); err != nil {
			return nil, err
		}

		snippets = append(snippets, snippet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
