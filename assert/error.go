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
		assert.fail(details, "unexpected error[%s]", err)
	}
}

// NoError will assert that given error is nil.
// If it's nil then the Fatal() function is called with details.
func NoError(t testing.TB, err error, details ...interface{}) {
	t.Helper()
	assert := New(t, Fatal)
	assert.NoError(err, details...)
}

// Error will call the failure function with the given details if the error is nil.
func (assert *Assert) Error(err error, details ...interface{}) {
	assert.t.Helper()
	if err == nil {
		assert.fail(details, "expected error, got nil")
	}
}

// Error will call Fatal with the given details if the error is nil.
func Error(t testing.TB, err error, details ...interface{}) {
	t.Helper()
	assert := New(t, Fatal)
	assert.Error(err, details...)
}

// IsError will assert if the given error matches the want error.
// It uses the errors.Is() function to check if the error wraps the wanted error.
// It calls the failure function if errors.Is() returns false.
func (assert *Assert) IsError(got, want error, details ...interface{}) {
	assert.t.Helper()
	if !errors.Is(got, want) {
		assert.fail(details, "got [%v] but wanted [%v]", got, want)
	}
}

// IsError will assert if the given error matches the want error.
// It uses the errors.Is() function to check if the error wraps the wanted error.
// It calls the Fatal() function if errors.Is() returns false.
func IsError(t testing.TB, got, want error, details ...interface{}) {
	t.Helper()
	assert := New(t, Fatal)
	assert.IsError(got, want, details...)
}
