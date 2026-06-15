package main

import (
	"net/http"
	"testing"

	"github.com/ItakawaM/snippetbox/internal/assert"
)

func Test_ping(t *testing.T) {
	app := newTestApplcation(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, string(body), "OK")
}
