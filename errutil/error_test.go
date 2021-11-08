package errutil_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/madlambda/spells/errutil"
)

func TestErrorSentinelWrapping(t *testing.T) {
	const (
		someError        errutil.Error = "someError"
		someAnotherError errutil.Error = "someAnotherError"
	)

	wrappedErr := fmt.Errorf("wrapping up: %w", someError)
	assertErrorIsWrapped(t, wrappedErr, someError)

	wrappedErr2 := fmt.Errorf("wrapping up: %w", someAnotherError)
	assertErrorIsWrapped(t, wrappedErr2, someAnotherError)
}

func TestErrorRepresentation(t *testing.T) {
	const (
		someError errutil.Error = "someError"
	)

	got := someError.Error()
	want := string(someError)

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func assertErrorIsWrapped(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Errorf("error [%v] is not wrapping [%v]", err, target)
	}
}
