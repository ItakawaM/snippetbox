package mocks

import "github.com/ItakawaM/snippetbox/internal/models"

type CommentModel struct{}

func (m *CommentModel) Latest(snippetID int) ([]*models.Comment, error) {
	return nil, nil
}
