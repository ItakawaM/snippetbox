package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
	Insert(name string, email string, password string) (int, error)
	Authenticate(email string, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (*User, error)
	PasswordUpdate(id int, currentPassword string, newPassword string) error
}

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

func (m *UserModel) Get(id int) (*User, error) {
	statement := `SELECT id, name, email, created FROM users WHERE id = $1;`

	user := &User{}
	if err := m.DB.QueryRow(statement, id).Scan(&user.ID, &user.Name, &user.Email, &user.Created); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}

		return nil, err
	}

	return user, nil
}

func (m *UserModel) Authenticate(email string, password string) (int, error) {
	statement := `SELECT id, hashed_password FROM users WHERE email = $1`

	var id int
	var hashedPassword []byte
	if err := m.DB.QueryRow(statement, strings.ToLower(email)).Scan(&id, &hashedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	statement := `SELECT EXISTS(SELECT true FROM users WHERE id = $1)`

	var exists bool
	err := m.DB.QueryRow(statement, id).Scan(&exists)

	return exists, err
}

func (m *UserModel) PasswordUpdate(id int, currentPassword string, newPassword string) error {
	statement := `SELECT hashed_password FROM users WHERE id = $1;`

	var hashedPassword []byte
	if err := m.DB.QueryRow(statement, id).Scan(&hashedPassword); err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(currentPassword)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		} else {
			return err
		}
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	statement = `UPDATE users SET hashed_password = $1 WHERE id = $2;`

	_, err = m.DB.Exec(statement, newHashedPassword, id)
	return err
}
