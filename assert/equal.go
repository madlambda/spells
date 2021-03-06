package assert

import (
	"fmt"
	"math"
	"testing"
)

var ε = math.Nextafter(1, 2) - 1

// EqualStrings compares the two strings for equality.
// If they are not equal t.Fatal is called using the details parameter.
// The details parameter can be a single string of a format string + parameters.
func EqualStrings(t *testing.T, want string, got string, details ...interface{}) {
	t.Helper()
	if want != got {
		detail := errordetails(details...)
		t.Fatalf("wanted[%s] but got[%s].%s", want, got, detail)
	}
}

// EqualInts compares the two ints for equality.
// If they are not equal t.Fatal is called using the details parameter.
// The details parameter can be a single string of a format string + parameters.
func EqualInts(t *testing.T, want int, got int, details ...interface{}) {
	t.Helper()
	if want != got {
		detail := errordetails(details...)
		t.Fatalf("wanted[%d] but got[%d].%s", want, got, detail)
	}
}

// EqualFloats compares the two floats for equality.
// If they are not equal t.Fatal is called using the details parameter.
// The details parameter can be a single string of a format string + parameters.
func EqualFloats(
	t *testing.T, want, got float64, details ...interface{},
) {
	t.Helper()

	if !floatEqual(want, got) {
		detail := errordetails(details...)
		t.Fatalf("wanted[%f] but got[%f].%s",
			want, got, detail)
	}
}

// EqualErrs compares if two errors have the same error description (by calling .Error()).
// If they are not equal t.Fatal is called using the details parameter.
// Both errors can't be nil.
// The details parameter can be a single string of a format string + parameters.
func EqualErrs(
	t *testing.T, want, got error, details ...interface{},
) {
	t.Helper()

	detail := errordetails(details...)
	if got != nil {
		if want != nil {
			if got.Error() != want.Error() {
				t.Fatalf("wanted[%s] but got[%s].%s", want,
					got, detail)
			}

			return
		}

		t.Fatalf("got unexpected error[%s].%s", got, detail)
		return
	}

	if want != nil {
		t.Fatalf("expected error[%s] but got nil.%s",
			want, detail)
	}
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

func floatEqual(a, b float64) bool {
	return math.Abs(a-b) < ε && math.Abs(b-a) < ε
}
