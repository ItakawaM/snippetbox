package models

import (
	"database/sql"
	"time"
)

type CommentModelInterface interface {
	Latest(snippetID int) ([]*Comment, error)
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
