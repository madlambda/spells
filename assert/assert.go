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
}

// FailureReport is the function type used to report assert errors.
// See assert.Fatal and assert.Err for implementations.
type FailureReport func(assert *Assert, message string)

const detailSeparator = ": "

// New creates a new assert helper object with a custom fail function and an
// optional context detail.
// For calling t.Fatal() or t.Error() in case of failures, see Fatal() and Err()
// respectively.
// Example:
//   assert := assert.New(t, assert.Fatal)
// The variadic details parameter must be a string format followed by its format
// arguments.
// Example:
//   v1.Name = "test"
//   v2.Name = "tesd"
//   assert := assert.New(t, assert.Err, "comparing objects %s and %s", v1, v2)
//   assert.EqualStrings(v1.Name, v2.Name, "Name mismatch")
//   ...
// The code above fails with message below:
//   wanted[test] but got[tesd].Name mismatch: comparing objects Value1 and Value2
func New(t *testing.T, fail FailureReport, details ...interface{}) *Assert {
	return &Assert{
		t:        t,
		failfunc: fail,
		details:  details,
	}
}

func (assert *Assert) fail(context []interface{}, details ...interface{}) {
	assert.t.Helper()
	assert.failfunc(assert, errctx(assert.details,
		errctx(context, details...)))
}

func (assert *Assert) failif(cond bool, context []interface{}, details ...interface{}) {
	if cond {
		assert.fail(context, details...)
	}
}

// Fatal is a FailureReport that calls t.Fatal() to abort the test.
func Fatal(assert *Assert, message string) {
	assert.t.Helper()
	assert.t.Fatal(message)
}

// Err is a FailureReport that calls t.Error().
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

func errctx(context []interface{}, details ...interface{}) string {
	errstr := errordetails(details...)
	if len(errstr) > 0 {
		errstr += detailSeparator
	}
	errstr += errordetails(context...)
	return errstr
}
