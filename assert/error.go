package assert

import (
	"errors"
	"testing"
)

// NoError will assert that given error is nil.
// If it's nil then the failure function is called with details.
func (assert *Assert) NoError(err error, details ...interface{}) {
	assert.t.Helper()
	if err != nil {
		assert.fail("unexpected error[%s].%s", err, errordetails(details...))
	}
}

// NoError will assert that given error is nil.
// If it's nil then the Fatal() function is called with details.
func NoError(t *testing.T, err error, details ...interface{}) {
	t.Helper()
	assert := New(t, Fatal)
	assert.NoError(err, details...)
}

// Error will call Fatal if the given error is nil.
// The details parameter can be a single string of a format string + parameters.
func Error(t *testing.T, err error, details ...interface{}) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error, got nil.%s", errordetails(details...))
	}
}

// IsError will assert if the given error matches the want error.
// It uses the errors.Is() function to check if the error wraps the wanted error.
// It calls the failure function if errors.Is() returns false.
func (assert *Assert) IsError(got, want error, details ...interface{}) {
	assert.t.Helper()
	if !errors.Is(got, want) {
		assert.fail("got [%v] but wanted [%v]: %s", got, want, errordetails(details...))
	}
}

// IsError will assert if the given error matches the want error.
// It uses the errors.Is() function to check if the error wraps the wanted error.
// It calls the Fatal() function if errors.Is() returns false.
func IsError(t *testing.T, got, want error, details ...interface{}) {
	t.Helper()
	assert := New(t, Fatal)
	assert.IsError(got, want, details...)
}
