package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ItakawaM/snippetbox/internal/assert"
)

func Test_secureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	secureHeaders(next).ServeHTTP(rr, r)
	rs := rr.Result()

	want := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), want)

	want = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), want)

	want = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), want)

	want = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), want)

	want = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), want)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
