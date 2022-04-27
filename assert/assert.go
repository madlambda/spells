package assert

import (
	"testing"
)

type Assert struct {
	t        *testing.T
	details  []interface{}
	failfunc FailureReport
	Failures int
}

type FailureReport func(assert *Assert, details ...interface{})

// New creates a new assert helper object with the context message constructed
// from the optional details slice.
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
	assert.failfunc(assert, details...)
}

func (assert *Assert) Bool(want bool, got bool, details ...interface{}) {
	assert.t.Helper()
	if want != got {
		assert.fail("want[%t] but got[%t].%s", want, got, errordetails(details...))
	}
}

func (assert *Assert) True(b bool, details ...interface{}) {
	assert.t.Helper()
	assert.Bool(true, b, details...)
}

func (assert *Assert) False(b bool, details ...interface{}) {
	assert.Bool(false, b, details...)
}

func (assert *Assert) Success() bool {
	return assert.Failures == 0
}

func True(t *testing.T, cond bool, details ...interface{}) {
	assert := New(t, Fatal)
	assert.True(cond, details...)
}

func Fatal(assert *Assert, details ...interface{}) {
	assert.t.Helper()
	assert.t.Fatalf("%s.%s", errordetails(details...), errordetails(assert.details...))
}

func Err(assert *Assert, details ...interface{}) {
	assert.t.Helper()
	assert.Failures++
	assert.t.Errorf("%s.%s", errordetails(details...), errordetails(assert.details...))
}
