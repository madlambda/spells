package assert

import (
	"math"
	"testing"
)

var ε = math.Nextafter(1, 2) - 1

// EqualBools asserts two booleans for equality.
// If they are not the equal then the failure function is called.
func (assert *Assert) EqualBools(want bool, got bool, details ...interface{}) {
	assert.t.Helper()
	if want != got {
		assert.fail(details, "want[%t] but got[%t]", want, got)
	}
}

// IsTrue asserts that b is true.
// If it's not then the failure function is called with details.
func (assert *Assert) IsTrue(b bool, details ...interface{}) {
	assert.t.Helper()
	assert.EqualBools(true, b, details...)
}

// IsFalse asserts that b is false.
// If it's not then the failure function is called with details.
func (assert *Assert) IsFalse(b bool, details ...interface{}) {
	assert.EqualBools(false, b, details...)
}

// IsTrue asserts that b is true.
// If it's not then Fatal() is called with details.
func IsTrue(t *testing.T, cond bool, details ...interface{}) {
	assert := New(t, Fatal)
	assert.IsTrue(cond, details...)
}

// EqualStrings compares the two strings for equality.
// If they are not equal then the failure function is called with details.
func (assert *Assert) EqualStrings(want string, got string, details ...interface{}) {
	assert.t.Helper()
	if want != got {
		assert.fail(details, "wanted[%s] but got[%s]", want, got)
	}
}

// EqualStrings compares the two strings for equality.
// If they are not equal then the Fatal() function is called with details.
func EqualStrings(t *testing.T, want string, got string, details ...interface{}) {
	t.Helper()
	assert := New(t, Fatal)
	assert.EqualStrings(want, got, details...)
}

// EqualInts compares the two ints for equality.
// If they are not equal then the failure function is called with details.
func (assert *Assert) EqualInts(want int, got int, details ...interface{}) {
	assert.t.Helper()
	if want != got {
		assert.fail(details, "wanted[%d] but got[%d]", want, got)
	}
}

// EqualInts compares the two ints for equality.
// If they are not equal then the Fatal() function is called with details.
func EqualInts(t *testing.T, want int, got int, details ...interface{}) {
	t.Helper()
	assert := New(t, Fatal)
	assert.EqualInts(want, got, details...)
}

// EqualUints compares the two uint64s for equality.
// If they are not equal then the failure function is called with details.
func (assert *Assert) EqualUints(want uint64, got uint64, details ...interface{}) {
	assert.t.Helper()
	if want != got {
		assert.fail(details, "wanted[%d] but got[%d]", want, got)
	}
}

// EqualFloats compares the two floats for equality.
// If they are not equal then the failure function is called with details.
func (assert *Assert) EqualFloats(want float64, got float64, details ...interface{}) {
	assert.t.Helper()
	if !floatEqual(want, got) {
		assert.fail(details, "wanted[%f] but got[%f]", want, got)
	}
}

// EqualFloats compares the two floats for equality.
// If they are not equal then the Fatal() function is called with details.
func EqualFloats(t *testing.T, want, got float64, details ...interface{}) {
	t.Helper()
	assert := New(t, Fatal)
	assert.EqualFloats(want, got)
}

// EqualComplexes compares the two complex numbers for equality.
// If they are not equal then the failure function is called with details.
func (assert *Assert) EqualComplexes(want, got complex128, details ...interface{}) {
	assert.t.Helper()
	if want != got {
		assert.fail(details, "wanted complex number [%d] but got [%d]", want, got)
	}
}

// EqualComplexes compares the two complex numbers for equality.
// If they are not equal then the Fatal function is called with details.
func EqualComplexes(t *testing.T, want, got complex128, details ...interface{}) {
	t.Helper()
	assert := New(t, Fatal)
	assert.EqualComplexes(want, got, details...)
}

// EqualErrs compares if two errors have the same error description (by calling .Error()).
// If they are not equal then the failure function is called with details.
// Both errors can't be nil.
func (assert *Assert) EqualErrs(want error, got error, details ...interface{}) {
	if got != nil {
		if want != nil {
			if got.Error() != want.Error() {
				assert.fail(details, "wanted[%s] but got[%s]", want, got)
			}

			return
		}

		assert.fail(details, "got unexpected error[%s].%s", got)
		return
	}

	if want != nil {
		assert.fail(details, "expected error[%s] but got nil", want)
	}
}

// EqualErrs compares if two errors have the same error description (by calling .Error()).
// If they are not equal then the Fatal() function is called with details.
// Both errors can't be nil.
func EqualErrs(t *testing.T, want, got error, details ...interface{}) {
	t.Helper()
	assert := New(t, Fatal)
	assert.EqualErrs(want, got, details...)
}

func floatEqual(a, b float64) bool {
	return math.Abs(a-b) < ε && math.Abs(b-a) < ε
}
