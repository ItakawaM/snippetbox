package mocks

import "github.com/ItakawaM/snippetbox/internal/models"

type CommentModel struct{}

func (m *CommentModel) Latest(snippetID int) ([]*models.Comment, error) {
	return nil, nil
}

func (m *CommentModel) Insert(ownerID int, snippetID int, content string) (int, error) {
	return 0, nil
}
