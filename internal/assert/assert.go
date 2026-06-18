package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, got T, want T) {
	t.Helper()

	if got != want {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func StringContains(t *testing.T, got string, want string) {
	t.Helper()

	if !strings.Contains(got, want) {
		t.Errorf("got: %v, want substring: %v", got, want)
	}
}

func NilError(t *testing.T, got error) {
	t.Helper()

	if got != nil {
		t.Errorf("got: %v, want nil", got)
	}
}
