package assert

import "testing"

func NoError(t *testing.T, err error, details ...interface{}) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error[%s].%s", err, errordetails(details...))
	}
}

func Error(t *testing.T, err error, details ...interface{}) {
	if err == nil {
		t.Fatalf("expected error, got nil.%s", errordetails(details...))
	}
}
