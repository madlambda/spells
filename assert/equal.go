package assert

import (
	"fmt"
	"testing"
)

func EqualStrings(t *testing.T, want string, got string, details ...interface{}) {
	t.Helper()
	if want != got {
		detail := errordetails(details...)
		t.Fatalf("wanted[%s] but got[%s].%s", want, got, detail)
	}
}

func EqualInts(t *testing.T, want int, got int, details ...interface{}) {
	t.Helper()
	if want != got {
		detail := errordetails(details...)
		t.Fatalf("wanted[%d] but got[%d].%s", want, got, detail)
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
