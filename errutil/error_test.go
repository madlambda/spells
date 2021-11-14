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

func TestErrorChain(t *testing.T) {
	errs := []error{
		errors.New("top error"),
		errors.New("wrapped error 1"),
		errors.New("wrapped error 2"),
		errors.New("wrapped error 3"),
	}

	err := errutil.Chain(errs...)
	if err == nil {
		t.Fatal("got nil, wanted error")
	}

	got := err
	for i, want := range errs {
		if got == nil {
			t.Fatal("expected error to exist, got nil")
		}

		gotChain := got.(errutil.ErrorChain)
		if gotChain.Head != want {
			t.Fatalf("error[%d] got: [%v] want: [%v]", i, gotChain.Head, want)
		}

		got = errors.Unwrap(got)
	}

	if got != nil {
		t.Fatalf("wanted error chain to reach end (nil), got chain [%v] instead", got)
	}
}

func assertErrorIsWrapped(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Errorf("error [%v] is not wrapping [%v]", err, target)
	}
}
