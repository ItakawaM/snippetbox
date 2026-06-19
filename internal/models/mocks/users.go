package mocks

import (
	"time"

	"github.com/ItakawaM/snippetbox/internal/models"
)

type UserModel struct{}

func (m *UserModel) Insert(name string, email string, password string) (int, error) {
	switch email {
	case "dupe@example.com":
		return 0, models.ErrDuplicateEmail
	default:
		return 1, nil
	}
}

func (m *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return &models.User{
			ID:      1,
			Name:    "Alice",
			Email:   "alice@example.com",
			Created: time.Now(),
		}, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *UserModel) Authenticate(email string, password string) (int, error) {
	if email == "correct@example.com" && password == "pa$$Word1" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}
