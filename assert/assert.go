package assert

import (
	"fmt"
	"testing"
)

// Assert is a custom assert helper.
type Assert struct {
	t        *testing.T
	details  []interface{}
	failfunc FailureReport
	Failures int
}

// FailureReport is the function type used to report assert errors.
// See assert.Fatal and assert.Err for implementations.
type FailureReport func(assert *Assert, message string)

// New creates a new assert helper object with a custom fail function and an
// optional context detail.
// For calling t.Fatal() or t.Error() in case of failures, see Fatal() and Err()
// respectively.
// Example:
//   assert := assert.New(t, assert.Fatal)
func New(t *testing.T, fail FailureReport, details ...interface{}) *Assert {
	return &Assert{
		t:        t,
		failfunc: fail,
		details:  details,
	}
}

func (assert *Assert) fail(details ...interface{}) {
	assert.t.Helper()
	assert.Failures++

	message := errordetails(details...)
	if len(assert.details) > 0 {
		message += fmt.Sprintf(": %s", errordetails(assert.details...))
	}

	assert.failfunc(assert, message)
}

// Success tells if there was no assertion failure.
func (assert *Assert) Success() bool {
	return assert.Failures == 0
}

func Fatal(assert *Assert, message string) {
	assert.t.Helper()
	assert.t.Fatal(message)
}

func Err(assert *Assert, message string) {
	assert.t.Helper()
	assert.t.Error(message)
}

func errordetails(details ...interface{}) string {
	if len(details) == 1 {
		return details[0].(string)
	}

	if len(details) > 1 {
		return fmt.Sprintf(details[0].(string), details[1:]...)
	}
	return ""
}
