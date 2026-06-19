package models

import (
	"database/sql"
	"errors"
	"time"
)

type CommentModelInterface interface {
	Latest(snippetID int) ([]*Comment, error)
	Insert(ownerID int, snippetID int, content string) (int, error)
}

type Comment struct {
	ID        int
	OwnerID   int
	OwnerName string
	SnippetID int
	Content   string
	Created   time.Time
}

type CommentModel struct {
	DB *sql.DB
}

func (m *CommentModel) Latest(snippetID int) ([]*Comment, error) {
	statement :=
		`SELECT u.name, c.id, c.created, c.content 
	FROM comments c 
	JOIN users u ON c.owner_id = u.id
	WHERE c.snippet_id = $1
	ORDER BY c.created DESC;`

	rows, err := m.DB.Query(statement, snippetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*Comment{}
	for rows.Next() {
		comment := &Comment{}
		if err := rows.Scan(&comment.OwnerName, &comment.ID, &comment.Created, &comment.Content); err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (m *CommentModel) Insert(ownerID int, snippetID int, content string) (int, error) {
	statement :=
		`INSERT INTO comments (owner_id, snippet_id, content, created)
    SELECT $1, $2, $3, NOW()
    FROM snippets
    WHERE id = $2 AND expires > NOW()
    RETURNING id;`

	var id int
	if err := m.DB.QueryRow(statement, ownerID, snippetID, content).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrSnippetExpired
		}

		return 0, err
	}

	return id, nil
}
