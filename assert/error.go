package assert

import "testing"

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
