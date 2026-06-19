package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/ItakawaM/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplcation(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, string(body), "OK")
}

func TestApplication_SnippetView(t *testing.T) {
	app := newTestApplcation(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.24",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/Golang",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)
			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}

func TestUserSignup(t *testing.T) {
	const urlPath = "/user/signup"

	app := newTestApplcation(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, urlPath)
	csrfToken := extractCSRFToken(t, body)

	const (
		validName     = "Bob"
		validPassword = "validPa!!Word123"
		validEmail    = "bob@example.com"
		formTag       = "<form action='/user/signup' method='POST' novalidate>"
	)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  string
	}{
		{
			name:         "Valid Signup",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusSeeOther,
		},
		{
			name:         "Invalid CSRF Token",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "invalid",
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Empty Name",
			userName:     "",
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty Email",
			userName:     validName,
			userEmail:    "",
			userPassword: validPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid Email",
			userName:     validName,
			userEmail:    "invalid",
			userPassword: validPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid Email Duplicate",
			userName:     validName,
			userEmail:    "dupe@example.com",
			userPassword: validPassword,
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty Password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "",
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid Password Short",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "123aS!",
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid Password No Uppercase",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "123456as!",
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid Password No Lowercase",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "123456AS!",
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid Password No Special",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "12345678AS",
			csrfToken:    csrfToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, urlPath, form)

			assert.Equal(t, code, tt.wantCode)
			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}
		})
	}
}

func TestApplication_SnippetCreate(t *testing.T) {
	app := newTestApplcation(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		code, header, _ := ts.get(t, "/snippet/create")

		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, header.Get("Location"), "/user/login")
	})

	t.Run("Authenticated", func(t *testing.T) {
		_, _, body := ts.get(t, "/user/login")
		csrf_token := extractCSRFToken(t, body)

		form := url.Values{}
		form.Add("email", "correct@example.com")
		form.Add("password", "pa$$Word1")
		form.Add("csrf_token", csrf_token)
		ts.postForm(t, "/user/login", form)

		code, _, body := ts.get(t, "/snippet/create")
		assert.Equal(t, code, http.StatusOK)
		assert.StringContains(t, body, `<form action="/snippet/create" method="POST">`)
	})
}
