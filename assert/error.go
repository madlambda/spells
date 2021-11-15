package assert

import (
	"errors"
	"testing"
)

// NoError will call Fatal if the given error is not nil.
// The details parameter can be a single string of a format string + parameters.
func NoError(t *testing.T, err error, details ...interface{}) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error[%s].%s", err, errordetails(details...))
	}
}

// Error will call Fatal if the given error is nil.
// The details parameter can be a single string of a format string + parameters.
func Error(t *testing.T, err error, details ...interface{}) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error, got nil.%s", errordetails(details...))
	}
}

// IsError will call Fatal if the given error does not match the want error.
// It uses the errors.Is() function to check if the error wraps the wanted error.
func IsError(t *testing.T, got, want error, details ...interface{}) {
	t.Helper()

	detail := errordetails(details...)
	if !errors.Is(got, want) {
		t.Fatalf("got [%v] but wanted [%v]: %s", got, want, detail)
	}
}
