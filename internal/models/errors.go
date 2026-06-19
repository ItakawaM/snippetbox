package models

import "errors"

var ErrNoRecord error = errors.New("models: no matching record found")

var ErrInvalidCredentials error = errors.New("models: invalid credentials")
var ErrDuplicateEmail error = errors.New("models: duplicate email")

var ErrSnippetExpired = errors.New("models: snippet is expired")
