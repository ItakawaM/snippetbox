package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name string, email string, password string) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	statement := `INSERT INTO users (name, email, hashed_password, created)
	VALUES ($1, $2, $3, NOW()) RETURNING id;`

	var id int
	if err := m.DB.QueryRow(statement, name, strings.ToLower(email), string(hashedPassword)).Scan(&id); err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_uc_email" {
				return 0, ErrDuplicateEmail
			}
		}

		return 0, err
	}

	return id, nil
}

func (m *UserModel) Authenticate(email string, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
