package mocks

import "github.com/ItakawaM/snippetbox/internal/models"

type UserModel struct{}

func (m *UserModel) Insert(name string, email string, password string) (int, error) {
	switch email {
	case "dupe@example.com":
		return 0, models.ErrDuplicateEmail
	default:
		return 1, nil
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
